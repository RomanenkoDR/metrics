package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {

	// Получаем SECRET_KEY из переменной окружения
	jwtKey := "SECRET_KEY"
	if jwtKey == "" {
		return "", fmt.Errorf("SECRET_KEY не установлен в переменные окружения")
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtKey))
}

func СheckAuthToken(r *http.Request) (string, int, error) {
	jwtToken := r.Header.Get("Authorization")

	if jwtToken == "" || !strings.HasPrefix(jwtToken, "Bearer ") {
		return "", http.StatusUnauthorized, errors.New("пользователь не авторизован")
	}

	jwtToken = strings.TrimPrefix(jwtToken, "Bearer ")

	t, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte("SECRET_KEY"), nil
	})

	switch {
	case t.Valid:
		validClaims := t.Claims.(jwt.MapClaims)
		username, ok := validClaims["username"].(string)
		if !ok || username == "" {
			return "", http.StatusBadRequest, errors.New("неудалось получить username из токена JWT")
		}
		return username, http.StatusOK, nil
	case errors.Is(err, jwt.ErrTokenMalformed):
		return "", http.StatusBadRequest, fmt.Errorf("токен имеет неправильную форму %w", err)
	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
		return "", http.StatusUnauthorized, fmt.Errorf("подпись токена недействительна %w", err)
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		return "", http.StatusUnauthorized, fmt.Errorf("срок действия токена истек или токен еще не действителен %w", err)
	default:
		return "", http.StatusInternalServerError, fmt.Errorf("неизвестная ошика на этапе проверки валидности токена: %w", err)
	}
}
