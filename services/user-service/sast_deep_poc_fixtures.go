// Package fixtures contains deliberate security vulnerability test fixtures for SAST evaluation.
// Generated for DevGuard automated supply chain security verification on Go Microservices.
package fixtures

import (
	"crypto/tls"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// ConfigureInsecureTLS demonstrates PoC 1: Insecure TLS MinVersion < 1.3
// CWE-326: Inadequate Encryption Strength
// CWE-327: Use of a Broken or Risky Cryptographic Algorithm
func ConfigureInsecureTLS() *tls.Config {
	return &tls.Config{
		MinVersion: tls.VersionTLS10, // VULNERABLE: Deprecated TLS 1.0 protocol
	}
}

// GetUserUnsafeRawSQL demonstrates PoC 2: SQL Injection via Unsanitized Raw Query
// CWE-89: Improper Neutralization of Special Elements used in an SQL Command ('SQL Injection')
func GetUserUnsafeRawSQL(db *sql.DB, username string) (*sql.Rows, error) {
	query := fmt.Sprintf("SELECT id, username, email FROM users WHERE username = '%s'", username)
	return db.Query(query) // VULNERABLE: Direct string formatting into SQL statement
}

// GenerateInsecureToken demonstrates PoC 3: Insecure JWT 'None' Algorithm & Hardcoded Secret
// CWE-347: Improper Verification of Cryptographic Signature
// CWE-798: Use of Hard-coded Credentials
func GenerateInsecureToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	// VULNERABLE: Insecure 'None' signing method allows forgery without signature verification
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	return token.SignedString(jwt.UnsafeAllowNoneSignatureType)
}

// ReadUserAvatarUnsafe demonstrates PoC 4: Path Traversal / Arbitrary File Access
// CWE-22: Improper Limitation of a Pathname to a Restricted Directory ('Path Traversal')
func ReadUserAvatarUnsafe(filename string) ([]byte, error) {
	baseDir := "/var/data/avatars"
	targetPath := filepath.Join(baseDir, filename)
	// VULNERABLE: Path traversal allows ../ to escape baseDir and access sensitive files
	return os.ReadFile(targetPath)
}

// SpawnUnboundedWorkers demonstrates PoC 5: Unbounded Goroutine Concurrency DoS
// CWE-400: Uncontrolled Resource Consumption ('Resource Exhaustion')
// CWE-362: Concurrent Execution using Shared Resource with Improper Synchronization
func SpawnUnboundedWorkers(requests []string) {
	for _, req := range requests {
		// VULNERABLE: Unbounded goroutine spawning without worker pool or semaphore control
		go func(payload string) {
			fmt.Println("Processing payload:", payload)
			time.Sleep(50 * time.Millisecond)
		}(req)
	}
}
