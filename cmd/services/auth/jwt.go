package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/vickon16/go-jwt-mysql/cmd/config"
	"github.com/vickon16/go-jwt-mysql/cmd/types"
	"github.com/vickon16/go-jwt-mysql/cmd/utils"
)

const UserKey string = "userId"

func CreateJWT(secret string, userID uuid.UUID) (string, error) {
	expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId":    userID.String(),
		"expiredAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get token from user header request
		tokenString := getTokenFromRequest(r)

		// validate the jwt
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("error validating token: %v", err)
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
			return
		}

		if !token.Valid {
			log.Printf("token not valid: %v", token)
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
			return
		}

		// get userId from token
		claims := token.Claims.(jwt.MapClaims)
		userId := uuid.MustParse(claims["userId"].(string))

		// fetch userId from the database
		user, err := store.GetUserByID(userId)
		if err != nil {
			log.Printf("error fetching user: %v", err)
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), UserKey, user.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func getTokenFromRequest(r *http.Request) string {
	tokenString := r.Header.Get("Authorization")
	if tokenString != "" {
		return tokenString
	}

	return ""
}

func validateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func GetUserIdFromContext(ctx context.Context) uuid.UUID {
	userId, ok := ctx.Value(UserKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}

	return userId
}
