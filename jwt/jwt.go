package jwthandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/shopeeProject/shopee/models"
	"github.com/shopeeProject/shopee/util"
)

var secretKey = []byte("your_secret_key")
var refreshTokens = map[string]string{}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(username string) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func GenerateRefreshToken(username string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   username,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 1 week
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

func ValidateRefreshToken(refresh_token string, r *util.Repository) util.DataResponse {
	refreshToken := []models.Token{}
	condition := models.Token{RefreshToken: refresh_token}
	err := r.DB.Where(condition).Find(&refreshToken).Error
	if err != nil {
		return util.DataResponse{Success: false, Message: "Error while finding the refresh token" + err.Error()}
	}
	if len(refreshToken) == 0 {
		return util.DataResponse{Success: false, Message: "Couldn't find the refresh token" + err.Error()}
	}
	return util.DataResponse{Success: true, Message: "Found the refresh Token", Data: map[string]string{"username": refreshToken[0].Email}}

}

func Refresh(refresh_token string, r *util.Repository) util.DataResponse {
	refreshTokenValidationResult := ValidateRefreshToken(refresh_token, r)
	if !refreshTokenValidationResult.Success {
		return util.DataResponse{Success: false, Message: refreshTokenValidationResult.Message}
	}
	newAccessToken, err := GenerateAccessToken(refreshTokenValidationResult.Data["username"])
	if err != nil {
		return util.DataResponse{Success: false, Message: "Error while generating new token" + err.Error()}
	}
	return util.DataResponse{Success: true, Message: "New Token Generated successfully", Data: map[string]string{"accessToken": newAccessToken}}
}

func Refresh1(w http.ResponseWriter, r *http.Request) {
	var body map[string]string
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshToken := body["refresh_token"]

	username, ok := refreshTokens[refreshToken]
	if !ok {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := GenerateAccessToken(username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": newAccessToken,
	})
}

func JwtMiddleware(tokenString string) util.Response {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil || !token.Valid {

		return util.Response{Success: false, Message: "Invalid Token"}
	}
	return util.Response{Success: true, Message: "Token Authentication Successful"}
}

func jwtMiddleware1(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
