package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtClaim struct {
	Id    int64  `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func CreateToken(id int64, email string, secret []byte) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JwtClaim{
		Id:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	})

	ss, err := token.SignedString(secret)
	if err != nil {
		panic(err)
	}

	return ss
}

// jwt claim token
func ClaimToken(token string, secret []byte) (*JwtClaim, error) {
	parsed, err := jwt.ParseWithClaims(token, &JwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	if parsed.Method != jwt.SigningMethodHS256 {
		return nil, errors.New("Invalid token")
	}

	if claim, ok := parsed.Claims.(*JwtClaim); ok {
		return claim, nil
	} else {
		return nil, errors.New("Invalid token")
	}
}
