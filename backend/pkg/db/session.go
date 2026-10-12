package db

import (
	"database/sql"
	"errors"
	"time"
)

const SessionLifetime = 30 * 24 * time.Hour

var ErrInvalidSession = errors.New("session expired or revoked")

func EnsureSessions() error {
	_, err := dbConn.Exec(`CREATE TABLE IF NOT EXISTS sessions (
 token_hash CHAR(64) CHARACTER SET ascii COLLATE ascii_bin PRIMARY KEY,
 user_id VARCHAR(36) NOT NULL,
 created_at DATETIME(6) NOT NULL,
 last_seen_at DATETIME(6) NOT NULL,
 expires_at DATETIME(6) NOT NULL,
 INDEX sessions_user (user_id), INDEX sessions_expiry (expires_at)
 ) ENGINE=InnoDB`)
	return err
}

func CreateSession(hash, userID string, expectedPassword ...string) error {
	// Serialize session creation with password changes/revocation for this user.
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var id, password string
	if err = tx.QueryRow("SELECT id,password FROM users WHERE id=? FOR UPDATE", userID).Scan(&id, &password); err != nil {
		return err
	}
	if len(expectedPassword) > 0 && password != expectedPassword[0] {
		return ErrInvalidSession
	}
	_, err = tx.Exec("INSERT INTO sessions(token_hash,user_id,created_at,last_seen_at,expires_at) VALUES(?,?,UTC_TIMESTAMP(6),UTC_TIMESTAMP(6),DATE_ADD(UTC_TIMESTAMP(6),INTERVAL 30 DAY))", hash, userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func DeleteExpiredSessions() error {
	_, err := dbConn.Exec("DELETE FROM sessions WHERE expires_at<=UTC_TIMESTAMP(6)")
	return err
}

func ValidateSession(hash string) (string, error) {
	tx, err := dbConn.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	var id string
	// Lock the session so expiry and renewal cannot race each other.
	err = tx.QueryRow("SELECT user_id FROM sessions WHERE token_hash=? AND expires_at>UTC_TIMESTAMP(6) FOR UPDATE", hash).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidSession
	}
	if err != nil {
		return "", err
	}
	_, err = tx.Exec("UPDATE sessions SET last_seen_at=UTC_TIMESTAMP(6),expires_at=DATE_ADD(UTC_TIMESTAMP(6),INTERVAL 30 DAY) WHERE token_hash=?", hash)
	if err != nil {
		return "", err
	}
	if err = tx.Commit(); err != nil {
		return "", err
	}
	return id, nil
}

// CheckSession does not count a background status poll as user activity.
func CheckSession(hash string) (string, error) {
	var id string
	err := dbConn.QueryRow("SELECT user_id FROM sessions WHERE token_hash=? AND expires_at>UTC_TIMESTAMP(6)", hash).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidSession
	}
	return id, err
}

func RevokeSession(hash string) error {
	_, err := dbConn.Exec("DELETE FROM sessions WHERE token_hash=?", hash)
	return err
}

func updatePasswordTx(tx *sql.Tx, id, hash string) error {
	if _, err := tx.Exec("UPDATE users SET password=? WHERE id=?", hash, id); err != nil {
		return err
	}
	_, err := tx.Exec("DELETE FROM sessions WHERE user_id=?", id)
	return err
}

// CompletePasswordReset consumes the reset token, changes the password and
// revokes every device's session as one atomic operation.
func CompletePasswordReset(secret, hash string) error {
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var id, userID string
	var valid, used bool
	err = tx.QueryRow("SELECT id,useridfs,expires_at>UTC_TIMESTAMP(),used FROM reset_password WHERE secret=? FOR UPDATE", secret).Scan(&id, &userID, &valid, &used)
	if errors.Is(err, sql.ErrNoRows) {
		return errors.New("invalid token")
	}
	if err != nil {
		return err
	}
	if used {
		return errors.New("token already used")
	}
	if !valid {
		return errors.New("token expired")
	}
	if err = updatePasswordTx(tx, userID, hash); err != nil {
		return err
	}
	if _, err = tx.Exec("UPDATE reset_password SET used=1 WHERE id=?", id); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdatePasswordWithExpected serializes verification with password resets and session creation.
func UpdatePasswordWithExpected(id, expected, replacement string) error {
	tx, err := dbConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var current string
	if err = tx.QueryRow("SELECT password FROM users WHERE id=? FOR UPDATE", id).Scan(&current); err != nil {
		return err
	}
	if current != expected {
		return ErrInvalidSession
	}
	if err = updatePasswordTx(tx, id, replacement); err != nil {
		return err
	}
	return tx.Commit()
}
