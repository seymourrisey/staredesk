package idgen

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// newID generates a prefixed ID using first 12 chars of UUID
// Example: USR-A1B2C3D4E5F6
func newID(prefix string) string {
	raw := strings.ReplaceAll(uuid.New().String(), "-", "")[:12]
	return fmt.Sprintf("%s-%s", prefix, strings.ToUpper(raw))
}

// newDateID generates a prefixed ID with date + 8 chars of UUID
// Example: TXN-20260307-A1B2C3D4
func newDateID(prefix string) string {
	date := time.Now().Format("20060102")
	raw := strings.ReplaceAll(uuid.New().String(), "-", "")[:8]
	return fmt.Sprintf("%s-%s-%s", prefix, date, strings.ToUpper(raw))
}

// NewUserID → USR-A1B2C3D4E5F6
func NewUserID() string { return newID("USR") }
