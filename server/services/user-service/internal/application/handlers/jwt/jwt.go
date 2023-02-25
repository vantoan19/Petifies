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
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	SessionID  uuid.UUID `json:"session_id"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpriresAt time.Time `json:"expired_at"`
}

func NewClaim(userID uuid.UUID, sessionId uuid.UUID, duration time.Duration) (*Claim, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Claim{
		ID:         tokenID,
		UserID:     userID,
		SessionID:  sessionId,
		IssuedAt:   time.Now(),
		ExpriresAt: time.Now().Add(duration),
	}, nil
}

func (p *Claim) Valid() error {
	if time.Now().After(p.ExpriresAt) {
		return ExpriedTokenError
	}

	return nil
}

type TokenMaker interface {
	CreateToken(userID uuid.UUID, sessionID uuid.UUID, duration time.Duration) (string, *Claim, error)
	VerifyToken(token string) (*Claim, error)
}

type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) TokenMaker {
	return &JWTMaker{secretKey: secretKey}
}

func (j *JWTMaker) CreateToken(userID uuid.UUID, sessionID uuid.UUID, duration time.Duration) (string, *Claim, error) {
	tokenClaim, err := NewClaim(userID, sessionID, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaim)
	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", nil, err
	}

	return signedToken, tokenClaim, nil
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
