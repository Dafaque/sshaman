package auth

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secretKey *ecdsa.PrivateKey
	issuer    string
}

func NewJWTManager(secretKey string, issuer string) (*JWTManager, error) {
	blk, _ := pem.Decode([]byte(secretKey))
	privateKeyBytes, err := x509.ParseECPrivateKey(blk.Bytes)
	if err != nil {
		return nil, err
	}
	if issuer == "" {
		issuer, _ = os.Hostname()
	}
	h := sha256.New()
	h.Write([]byte(issuer))
	return &JWTManager{secretKey: privateKeyBytes, issuer: hex.EncodeToString(h.Sum(nil))}, nil
}

func (j *JWTManager) GenerateToken(userId int64) (string, error) {
	// Token generation logic using secretKey
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodES384, jwt.MapClaims{
		"iss": j.issuer,
		"aud": j.issuer,
		"sub": strconv.FormatInt(userId, 10),
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"exp": now.Add(time.Hour * 24).Unix(),
	})
	key, err := j.privateKeyFunc(token)
	if err != nil {
		return "", err
	}
	return token.SignedString(key)
}

type authParams struct {
	userID int64
	valid  bool
}

func (j *JWTManager) ValidateToken(token string) (*authParams, error) {
	parser := jwt.NewParser(
		jwt.WithAudience(j.issuer),
		jwt.WithIssuer(j.issuer),
		jwt.WithValidMethods([]string{jwt.SigningMethodES384.Alg()}),
	)
	t, err := parser.Parse(token, j.publicKeyFunc)
	if err != nil {
		return nil, err
	}
	uid, err := t.Claims.GetSubject()
	if err != nil {
		return nil, err
	}
	userId, err := strconv.ParseInt(uid, 10, 64)
	if err != nil {
		return nil, err
	}
	return &authParams{userID: userId, valid: t.Valid}, nil
}

func (j *JWTManager) privateKeyFunc(_ *jwt.Token) (interface{}, error) {
	return j.secretKey, nil
}

func (j *JWTManager) publicKeyFunc(_ *jwt.Token) (interface{}, error) {
	return &j.secretKey.PublicKey, nil
}

// token signature is invalid: key is of invalid type: HMAC verify expects []byte
