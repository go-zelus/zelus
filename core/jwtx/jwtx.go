package jwtx

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = "HcVmUCTlkJ"

func SetSigningKey(key string) {
	signingKey = key
}

// New 创建Token
func New(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(signingKey))
}

// Parser 解析Token，claims必须传指针
func Parser(tokenString string, claims jwt.Claims) (jwt.Claims, error) {
	tk, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(signingKey), nil
	})

	if err != nil {
		if v, ok := err.(*jwt.ValidationError); ok {
			if v.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, fmt.Errorf("token malformed")
			} else if v.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, fmt.Errorf("token expired")
			} else if v.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, fmt.Errorf("token not valid yet")
			} else {
				return nil, fmt.Errorf("token error %v", v)
			}
		}
	}
	if tk != nil && tk.Valid {
		return tk.Claims, nil
	}
	return nil, fmt.Errorf("token invalid")
}
