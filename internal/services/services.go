package services

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/Lixerus/auth-service-task/internal/database"
	"github.com/Lixerus/auth-service-task/internal/middleware"
	"github.com/Lixerus/auth-service-task/internal/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func createAndSetTokenPair(c *gin.Context, UUID string, clientIP string) (string, string, error) {
	accessString, err := generateAccessToken(UUID, clientIP)
	if err != nil {
		return "", "", &TokenServiceError{Detail: "Failed to create access token", ErrMsg: err.Error()}
	}
	refreshTokenBytes := generateRefreshToken(clientIP)
	refreshStringBase64 := base64.StdEncoding.EncodeToString(refreshTokenBytes)

	hash, err := bcrypt.GenerateFromPassword(refreshTokenBytes, 10)
	if err != nil {
		return "", "", &TokenServiceError{Detail: "Failed to hash uuid", ErrMsg: err.Error()}
	}
	partialToken := getPartialAccessString(accessString)
	user := models.UserCredentials{ID: UUID, RefreshToken: string(hash), PartialAccessToken: partialToken}
	database.DB.Save(&user)
	setCookiesWithoutQueryEscape(c, "Authorization", accessString, int(time.Minute*15), "", "", false, true)
	setCookiesWithoutQueryEscape(c, "Authrefresh", refreshStringBase64, int(time.Hour*24), "", "", false, true)
	return accessString, refreshStringBase64, nil
}

func setCookiesWithoutQueryEscape(c *gin.Context, name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: http.SameSiteLaxMode,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func GetCredentials(c *gin.Context) {
	var requestIp string = c.ClientIP()
	query := c.Request.URL.Query()
	UUID := query["id"][0]
	if UUID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Failed to find id param"})
		return
	}
	accessString, refreshStringBase64, err := createAndSetTokenPair(c, UUID, requestIp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"access_token": accessString, "refresh_token": refreshStringBase64})
}

func RefreshAuthToken(c *gin.Context) {
	accessCookieString := c.GetString(middleware.AccessCookie)
	refreshCookieStringB64 := c.GetString(middleware.RefreshCookie)

	partialAccessString := getPartialAccessString(accessCookieString)
	var userCredential models.UserCredentials
	database.DB.Where(&models.UserCredentials{PartialAccessToken: partialAccessString}).First(&userCredential)
	if userCredential.ID == "" {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Invalid access cookie. Try to login instead"})
		return
	}
	refreshCookieBytes, err := base64.StdEncoding.DecodeString(refreshCookieStringB64)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Invalid refresh cookie. Bad encoding. Try to login instead."})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(userCredential.RefreshToken), refreshCookieBytes)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"detail": "Invalid refresh cookie. Try to login instead."})
		return
	}
	IP := parseRefreshToken(refreshCookieBytes)
	if IP != c.ClientIP() {
		go SendEmail()
		c.JSON(http.StatusForbidden, gin.H{"detail": "IP changed."})
		return
	}
	accessString, refreshStringBase64, err := createAndSetTokenPair(c, userCredential.ID, IP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"access_token": accessString, "refresh_token": refreshStringBase64})
}
