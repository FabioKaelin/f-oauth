// Package session implements revocable opaque sessions. Raw tokens are never stored.
package session

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"time"

	"github.com/fabiokaelin/f-oauth/config"
	"github.com/fabiokaelin/f-oauth/pkg/db"
)

const MaxAge = 30 * 24 * 60 * 60

func Hash(token string) string { h := sha256.Sum256([]byte(token)); return hex.EncodeToString(h[:]) }
func NewToken() (string, error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(raw), nil
}
func Create(userID string, expectedPassword ...string) (string, error) {
	token, err := NewToken()
	if err != nil {
		return "", err
	}
	if err = db.CreateSession(Hash(token), userID, expectedPassword...); err != nil {
		return "", err
	}
	return token, nil
}
func Validate(token string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(raw) != 32 || base64.RawURLEncoding.EncodeToString(raw) != token {
		return "", db.ErrInvalidSession
	}
	return db.ValidateSession(Hash(token))
}

func Check(token string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil || len(raw) != 32 || base64.RawURLEncoding.EncodeToString(raw) != token {
		return "", db.ErrInvalidSession
	}
	return db.CheckSession(Hash(token))
}
func Revoke(token string) error {
	if token == "" {
		return errors.New("missing session")
	}
	return db.RevokeSession(Hash(token))
}
func SetCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{Name: "token", Value: token, Path: "/", Domain: config.TokenURL, MaxAge: MaxAge, Expires: time.Now().UTC().Add(db.SessionLifetime), Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
}
func ClearCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{Name: "token", Value: "", Path: "/", Domain: config.TokenURL, MaxAge: -1, Expires: time.Unix(1, 0), Secure: true, HttpOnly: true, SameSite: http.SameSiteLaxMode})
}
