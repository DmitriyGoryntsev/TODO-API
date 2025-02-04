package helpers

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type SignedDetails struct {
	ID       int
	Username string
	Email    string
	Role     string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

func GenerateAllTokens(id int, username string, email string, role string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		ID:       id,
		Username: username,
		Email:    email,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic("cannot create token:", err)
		return
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Panic("cannot create refresh token:", err)
		return
	}

	return token, refreshToken, nil
}

func ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err != nil {
		log.Println("invalid token")
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		log.Println("invalid token")
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		log.Println("token expired")
		return nil, err
	}

	return claims, nil
}
