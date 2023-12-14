package auth

import (
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// func GenerateToken(secret string, id int, expTimeInSeconds *int) (*jwt.Token, error) {
func GenerateToken(secret string, id int, expTimeInSeconds *int) (string, error) {
	expTimeVal := 0
	switch expTimeInSeconds {
	case nil:
		expTimeVal = 24 * 3600
	default:
		if expTimeVal > 24*3600 {
			expTimeVal = 24 * 3600
		} else {
			expTimeVal = *expTimeInSeconds
		}
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expTimeVal) * time.Second)),
		Subject:   strconv.Itoa(id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(secret))
	if err != nil {
		// log.Println(err)
		return "", err
	}

	// log.Println("success")
	return ss, nil

}

func ValidateToken(tokenStr, secret string) (int, error) {
	type RegClaimsStruct struct {
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenStr, &RegClaimsStruct{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*RegClaimsStruct)
	if !ok {
		return 0, errors.New("unknown claim type")
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// func generateTokenWithClaims() {}
