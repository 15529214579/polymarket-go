package order

import "testing"

func TestLoadProxiesRequiresExplicitEnvironment(t *testing.T) {
	t.Setenv("CLOB_PROXY", "")
	t.Setenv("HTTPS_PROXY", "")
	if got := loadProxies(); len(got) != 0 {
		t.Fatalf("unexpected implicit proxies: %v", got)
	}
}

func TestLoadProxiesExplicitDirect(t *testing.T) {
	t.Setenv("CLOB_PROXY", "direct")
	t.Setenv("HTTPS_PROXY", "http://proxy.example:8080")
	if got := loadProxies(); len(got) != 0 {
		t.Fatalf("CLOB_PROXY=direct returned proxies: %v", got)
	}
}

func TestLoadProxiesExplicitURL(t *testing.T) {
	t.Setenv("CLOB_PROXY", "http://proxy.example:8080")
	t.Setenv("HTTPS_PROXY", "")
	got := loadProxies()
	if len(got) != 1 || got[0].String() != "http://proxy.example:8080" {
		t.Fatalf("explicit proxy = %v", got)
	}
}
