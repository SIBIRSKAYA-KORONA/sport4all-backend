package models

// swagger:model Session
type Session struct {
	// example: sdafasdjflasdjl
	SID string

	// example: 10
	ID uint

	// example: 123213
	ExpiresSec uint
}
