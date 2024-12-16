package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Lixerus/auth-service-task/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

type AccessTokenClaims struct {
	IpAddress string `json:"ip_address"`
	jwt.RegisteredClaims
}

type TokenServiceError struct {
	Detail string
	ErrMsg string
}

func (e *TokenServiceError) Error() string {
	return fmt.Sprintf("%s with error %s", e.Detail, e.ErrMsg)
}

func parseRefreshToken(token []byte) string {
	IP := token[25:]
	return string(IP)
}

func getPartialAccessString(token string) string {
	runes := []rune(token)
	return string(runes[len(runes)-10:])
}

func generateRefreshToken(IP string) []byte {
	tokenBytes := make([]byte, 25, 40)
	for i := range tokenBytes {
		tokenBytes[i] = byte(rand.Intn(128))
	}
	tokenBytes = append(tokenBytes, []byte(IP)...)
	return tokenBytes
}

func generateAccessToken(UUID string, requestIp string) (string, error) {
	tokenAccess := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub":        UUID,
		"ip_address": requestIp,
		"exp":        time.Now().Add(time.Minute * 15).Unix(),
	})
	accessString, err := tokenAccess.SignedString(config.DBConfig.ACCESS_SECRET.Text)
	return accessString, err
}
