package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	ExpriedTokenError = errors.New("token has expired")
)

type Claim struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	IssuedAt   time.Time
	ExpriredAt time.Time
}

func NewClaim(userID uuid.UUID, duration time.Duration) (*Claim, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Claim{
		ID:         tokenID,
		UserID:     userID,
		IssuedAt:   time.Now(),
		ExpriredAt: time.Now().Add(duration),
	}, nil
}

func (p *Claim) Valid() error {
	if time.Now().After(p.ExpriredAt) {
		return errors.New("token has expired")
	}

	return nil
}

type TokenMaker interface {
	CreateToken(userID uuid.UUID, duration time.Duration) (string, error)
	VerifyToken(token string) (*Claim, error)
}

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) TokenMaker {
	return &JWTMaker{secretKey: secretKey}
}

func (j *JWTMaker) CreateToken(userID uuid.UUID, duration time.Duration) (string, error) {
	tokenClaim, err := NewClaim(userID, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)
	return token.SignedString([]byte(j.secretKey))
}

func (j *JWTMaker) VerifyToken(token string) (*Claim, error) {
	keyFunc := func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token, the algorithm of the token doesn't match with the signing method")
		}
		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Claim{}, keyFunc)
	if err != nil {
		err_, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(err_.Inner, ExpriedTokenError) {
			return nil, ExpriedTokenError
		}
		return nil, fmt.Errorf("invalid token, %v", err)
	}

	claim, ok := jwtToken.Claims.(*Claim)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claim, nil
}
