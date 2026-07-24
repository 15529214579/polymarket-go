package sanitize

import (
	"strings"
	"testing"
)

func TestSecretString_RedactsTelegramBotToken(t *testing.T) {
	in := `Post "https://api.telegram.org/bot1234567890:ABCDEFGHIJKLMNOPQRSTUVWXYZabcdef/sendMessage": EOF`
	got := SecretString(in)
	if got == in {
		t.Fatal("expected redaction")
	}
	if want := "bot<redacted>"; !strings.Contains(got, want) {
		t.Fatalf("missing %q in %q", want, got)
	}
	if strings.Contains(got, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdef") {
		t.Fatalf("token was not redacted: %q", got)
	}
}

func TestSecretString_RedactsAPIKeyQuery(t *testing.T) {
	in := `Get "https://api.oddspapi.io/v4/fixtures?apiKey=abc-123&sportId=10": timeout`
	got := SecretString(in)
	if got == in {
		t.Fatal("expected redaction")
	}
	if strings.Contains(got, "abc-123") {
		t.Fatalf("api key was not redacted: %q", got)
	}
	if !strings.Contains(got, "apiKey=<redacted>") {
		t.Fatalf("missing redacted apiKey: %q", got)
	}
}
