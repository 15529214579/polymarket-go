package order

import "testing"

func TestNewSaltIsPositiveAndUnique(t *testing.T) {
	seen := make(map[string]struct{}, 1000)
	for i := 0; i < 1000; i++ {
		salt := NewSalt()
		if salt.Sign() <= 0 || !salt.IsInt64() {
			t.Fatalf("invalid salt %v", salt)
		}
		key := salt.String()
		if _, exists := seen[key]; exists {
			t.Fatalf("duplicate salt %s", key)
		}
		seen[key] = struct{}{}
	}
}
