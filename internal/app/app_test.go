package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Lixerus/auth-service-task/internal/database"
	"github.com/Lixerus/auth-service-task/internal/models"
)

func TestCredentials(t *testing.T) {
	InitDeps()
	router := InitApp()
	recorder := httptest.NewRecorder()
	testId := "testID1"
	router.ServeHTTP(recorder, httptest.NewRequest("POST", fmt.Sprintf("/credentials?id=%s", testId), nil))
	t.Run("Returns 200 status code", func(t *testing.T) {
		if recorder.Code != 201 {
			t.Error("Expected 201, got ", recorder.Code)
		}
	})
	t.Run("Get authorization and refresh tokens in cookies", func(t *testing.T) {
		result := recorder.Result()
		expectedCookiesMap := map[string]bool{"Authorization": false, "Authrefresh": false}
		cookies := result.Cookies()
		for _, v := range cookies {
			expectedCookiesMap[v.Name] = true
		}
		for k, v := range expectedCookiesMap {
			if !v {
				t.Errorf("No Cookie with name %s", k)
			}
		}
	})

	database.DB.Where("id = ?", testId).Delete(&models.UserCredentials{})
	db, _ := database.DB.DB()
	db.Close()
}

func setCookiePairOnTestRequest(r *http.Request, cookies []*http.Cookie) {
	cookieString := fmt.Sprintf("%s=%s; %s=%s", cookies[0].Name, cookies[0].Value, cookies[1].Name, cookies[1].Value)
	r.Header.Set("Cookie", cookieString)
}

func TestRefreshTokens(t *testing.T) {
	InitDeps()
	router := InitApp()

	recorder := httptest.NewRecorder()
	testId := "testID1"
	router.ServeHTTP(recorder, httptest.NewRequest("POST", fmt.Sprintf("/credentials?id=%s", testId), nil))
	t.Run("Get token pair success", func(t *testing.T) {
		if recorder.Code != 201 {
			t.Error("Expected 201, got ", recorder.Code)
		}
	})
	cookies := recorder.Result().Cookies()
	time.Sleep(1 * time.Second) // could get the same auth jwt token otherwise

	recorder = httptest.NewRecorder()
	request := httptest.NewRequest("POST", "/refresh", nil)
	setCookiePairOnTestRequest(request, cookies)
	router.ServeHTTP(recorder, request)
	refreshedCookies := recorder.Result().Cookies()
	t.Run("Refresh handler correct status", func(t *testing.T) {
		if recorder.Code != 201 {
			t.Error("Expected 201, got ", recorder.Code)
		}
	})
	t.Run("Refreshed old tokens with new pair", func(t *testing.T) {
		oldCookiesMap := map[string]string{cookies[0].Name: cookies[0].Value, cookies[1].Name: cookies[1].Value}
		if len(refreshedCookies) < 2 {
			t.Error("Not enough tokens sent")
		}
		for _, v := range refreshedCookies {
			if oldCookiesMap[v.Name] == "" {
				t.Error("Missing cookie with name ", v.Name)
			}
			if oldCookiesMap[v.Name] == v.Value {
				t.Error("Expected new cookie. Got same cookie with name", v.Name)
			}
		}
	})

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/refresh", nil)
	setCookiePairOnTestRequest(request, cookies)
	router.ServeHTTP(recorder, request)
	t.Run("Cannot use old token pair", func(t *testing.T) {
		if recorder.Code != 403 {
			t.Error("Expected 403, got ", recorder.Code)
		}
	})

	recorder = httptest.NewRecorder()
	testId2 := "testID2"
	router.ServeHTTP(recorder, httptest.NewRequest("POST", fmt.Sprintf("/credentials?id=%s", testId2), nil))
	newCookies := recorder.Result().Cookies()

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/refresh", nil)
	mixedCookiePair := []*http.Cookie{refreshedCookies[0], newCookies[1]}
	setCookiePairOnTestRequest(request, mixedCookiePair)
	router.ServeHTTP(recorder, request)
	t.Run("Cannot use 2 valid tokens from different paris", func(t *testing.T) {
		if recorder.Code != 403 {
			t.Error("Expected 403, got ", recorder.Code)
		}
	})

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/refresh", nil)
	tamperedRefreshCookie := []byte(newCookies[1].Value)
	tamperedRefreshCookie[10] = 51
	newCookies[1].Value = string(tamperedRefreshCookie)
	setCookiePairOnTestRequest(request, newCookies)
	router.ServeHTTP(recorder, request)
	t.Run("Client cannot tamper the refresh token", func(t *testing.T) {
		if recorder.Code != 403 {
			t.Error("Expected 403, got ", recorder.Code)
		}
	})

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/refresh", nil)
	setCookiePairOnTestRequest(request, refreshedCookies)
	request.RemoteAddr = "172.0.1.0:345"
	router.ServeHTTP(recorder, request)
	t.Run("Client cannot change ip", func(t *testing.T) {
		if recorder.Code != 403 {
			t.Error("Expected 403, got ", recorder.Code)
		}
		responseBody := struct {
			Detail string `json:"detail"`
		}{}
		err := json.NewDecoder(recorder.Body).Decode(&responseBody)
		if err != nil {
			t.Error("Unexpected change of body response ", err.Error())
		}
		if responseBody.Detail != "IP changed." {
			t.Error("Unexpected response detail ", responseBody.Detail)
		}
	})

	recorder = httptest.NewRecorder()
	request = httptest.NewRequest("POST", "/refresh", nil)
	setCookiePairOnTestRequest(request, refreshedCookies)
	router.ServeHTTP(recorder, request)
	t.Run("Refreshed cookies are valid to use", func(t *testing.T) {
		if recorder.Code != 201 {
			t.Error("Expected 201, got ", recorder.Code)
		}
	})
	//cleanup
	database.DB.Where("id = ?", testId).Delete(&models.UserCredentials{})
	database.DB.Where("id = ?", testId2).Delete(&models.UserCredentials{})
	db, _ := database.DB.DB()
	db.Close()
}

func TestRefreshNoTokens(t *testing.T) {
	router := InitApp()
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest("POST", "/refresh", nil))
	t.Run("Refresh without Tokens in cookies", func(t *testing.T) {
		if recorder.Code != 400 {
			t.Error("Expected 400, got ", recorder.Code)
		}
	})
}
