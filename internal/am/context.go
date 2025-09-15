package am

// contextKey is a private type used for context keys to avoid collisions.
type contextKey string

const (
	EncryptionKeyCtxKey contextKey = "encryptionKey"
)
