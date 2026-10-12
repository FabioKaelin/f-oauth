package db

import (
	"errors"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/jmoiron/sqlx"
)

func TestSessionsWithMariaDB(t *testing.T) {
	dsn := os.Getenv("OAUTH_SESSION_TEST_DSN")
	if dsn == "" {
		t.Skip("set OAUTH_SESSION_TEST_DSN to an isolated test database")
	}
	conn, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	dbConn = conn
	defer func() { dbConn = nil }()
	var name string
	if err = conn.QueryRow("SELECT DATABASE()").Scan(&name); err != nil || name != "oauth_session_test" {
		t.Fatal("tests require the isolated oauth_session_test database")
	}
	conn.MustExec("DROP TABLE IF EXISTS sessions,reset_password,users")
	conn.MustExec("CREATE TABLE users(id VARCHAR(36) PRIMARY KEY,password TEXT NOT NULL)")
	conn.MustExec("CREATE TABLE reset_password(id VARCHAR(36) PRIMARY KEY,useridfs VARCHAR(36),secret VARCHAR(100),expires_at DATETIME,used BOOL)")
	if err = EnsureSessions(); err != nil {
		t.Fatal(err)
	}
	conn.MustExec("INSERT INTO users VALUES('owner','old-hash')")
	h := func(c string) string { return strings.Repeat(c, 64) }
	t.Run("creation and activity renew expiry", func(t *testing.T) {
		if err = CreateSession(h("a"), "owner"); err != nil {
			t.Fatal(err)
		}
		conn.MustExec("UPDATE sessions SET expires_at=DATE_ADD(UTC_TIMESTAMP(6),INTERVAL 1 DAY),last_seen_at=DATE_SUB(UTC_TIMESTAMP(6),INTERVAL 29 DAY) WHERE token_hash=?", h("a"))
		id, e := ValidateSession(h("a"))
		if e != nil || id != "owner" {
			t.Fatalf("validation: %s %v", id, e)
		}
		var remaining float64
		if e = conn.QueryRow("SELECT TIMESTAMPDIFF(SECOND,UTC_TIMESTAMP(),expires_at) FROM sessions WHERE token_hash=?", h("a")).Scan(&remaining); e != nil {
			t.Fatal(e)
		}
		if remaining < SessionLifetime.Seconds()-10 || remaining > SessionLifetime.Seconds()+10 {
			t.Fatalf("remaining expiry %v", remaining)
		}
	})
	t.Run("expired session cannot revive", func(t *testing.T) {
		conn.MustExec("UPDATE sessions SET expires_at=DATE_SUB(UTC_TIMESTAMP(6),INTERVAL 1 SECOND) WHERE token_hash=?", h("a"))
		if _, e := ValidateSession(h("a")); !errors.Is(e, ErrInvalidSession) {
			t.Fatalf("expired session accepted: %v", e)
		}
	})
	t.Run("parallel activity is safe", func(t *testing.T) {
		if e := CreateSession(h("b"), "owner"); e != nil {
			t.Fatal(e)
		}
		var wg sync.WaitGroup
		for i := 0; i < 12; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if _, e := ValidateSession(h("b")); e != nil {
					t.Error(e)
				}
			}()
		}
		wg.Wait()
	})
	t.Run("logout revokes", func(t *testing.T) {
		if e := RevokeSession(h("b")); e != nil {
			t.Fatal(e)
		}
		if _, e := ValidateSession(h("b")); !errors.Is(e, ErrInvalidSession) {
			t.Fatal("revoked token still accepted")
		}
	})
	t.Run("reset revokes all devices atomically", func(t *testing.T) {
		for _, c := range []string{"c", "d"} {
			if e := CreateSession(h(c), "owner"); e != nil {
				t.Fatal(e)
			}
		}
		conn.MustExec("INSERT INTO reset_password VALUES('reset','owner','secret',DATE_ADD(UTC_TIMESTAMP(),INTERVAL 1 HOUR),0)")
		if e := CompletePasswordReset("secret", "new-hash"); e != nil {
			t.Fatal(e)
		}
		for _, c := range []string{"c", "d"} {
			if _, e := ValidateSession(h(c)); !errors.Is(e, ErrInvalidSession) {
				t.Fatal("session survived reset")
			}
		}
		var password string
		conn.QueryRow("SELECT password FROM users WHERE id='owner'").Scan(&password)
		if password != "new-hash" {
			t.Fatal("password not updated")
		}
		if e := CompletePasswordReset("secret", "another-hash"); e == nil || e.Error() != "token already used" {
			t.Fatalf("reset token reuse: %v", e)
		}
	})
	t.Run("password change revokes and racing old login fails", func(t *testing.T) {
		if e := CreateSession(h("e"), "owner"); e != nil {
			t.Fatal(e)
		}
		if e := UpdatePassword("owner", "changed-hash"); e != nil {
			t.Fatal(e)
		}
		if _, e := ValidateSession(h("e")); !errors.Is(e, ErrInvalidSession) {
			t.Fatal("session survived change")
		}
		if e := CreateSession(h("f"), "owner", "new-hash"); !errors.Is(e, ErrInvalidSession) {
			t.Fatal("old password login raced reset")
		}
	})
	t.Run("verified password change rejects racing reset and revokes sessions", func(t *testing.T) {
		var before string
		conn.QueryRow("SELECT password FROM users WHERE id='owner'").Scan(&before)
		if e := UpdatePasswordWithExpected("owner", "stale-hash", "should-not-write"); !errors.Is(e, ErrInvalidSession) {
			t.Fatalf("racing reset: %v", e)
		}
		var after string
		conn.QueryRow("SELECT password FROM users WHERE id='owner'").Scan(&after)
		if before != after {
			t.Fatal("stale verification changed password")
		}
		if e := CreateSession(h("v"), "owner"); e != nil {
			t.Fatal(e)
		}
		if e := UpdatePasswordWithExpected("owner", before, "verified-new-hash"); e != nil {
			t.Fatal(e)
		}
		if _, e := CheckSession(h("v")); !errors.Is(e, ErrInvalidSession) {
			t.Fatal("password change did not revoke session")
		}
	})

	t.Run("expired reset preserves password and sessions", func(t *testing.T) {
		if e := CreateSession(h("g"), "owner"); e != nil {
			t.Fatal(e)
		}
		conn.MustExec("INSERT INTO reset_password VALUES('expired','owner','expired-secret',DATE_SUB(UTC_TIMESTAMP(),INTERVAL 1 SECOND),0)")
		if e := CompletePasswordReset("expired-secret", "bad-hash"); e == nil || e.Error() != "token expired" {
			t.Fatal("expired reset accepted")
		}
		if _, e := ValidateSession(h("g")); e != nil {
			t.Fatal("failed reset revoked session")
		}
	})
}
