package authorization

import (
	"errors"
	"time"

	"menu_admin/model"

	"github.com/dgrijalva/jwt-go"
)

var (
	ErrInvalidToken   = errors.New("invalid token")
	ErrCannotGetClaim = errors.New("can not get claim")
)

// GenerateToken .
func GenerateToken(data *model.User) (string, error) {
	claim := model.Claim{
		UserID: data.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 4).Unix(),
			Issuer:    "TF",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken .
func ValidateToken(t string) (model.Claim, error) {
	token, err := jwt.ParseWithClaims(t, &model.Claim{}, verifyFunction)
	if err != nil {
		return model.Claim{}, err
	}
	if !token.Valid {
		return model.Claim{}, ErrInvalidToken
	}

	claim, ok := token.Claims.(*model.Claim)
	if !ok {
		return model.Claim{}, ErrCannotGetClaim
	}

	return *claim, nil
}

func verifyFunction(t *jwt.Token) (interface{}, error) {
	return verifyKey, nil
}

func GenerateCodeVerification(email string) (string, error) {
	claim := model.CodeVerification{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			Issuer:    "TF",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	signedToken, err := token.SignedString(signKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateToken .
func ValidateCodeVerification(t string) (model.CodeVerification, error) {
	token, err := jwt.ParseWithClaims(t, &model.CodeVerification{}, verifyFunction)
	if err != nil {
		return model.CodeVerification{}, err
	}
	if !token.Valid {
		return model.CodeVerification{}, ErrInvalidToken
	}

	claim, ok := token.Claims.(*model.CodeVerification)
	if !ok {
		return model.CodeVerification{}, ErrCannotGetClaim
	}

	return *claim, nil
}
