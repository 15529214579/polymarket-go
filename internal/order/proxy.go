package order

import (
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	proxyList   []*url.URL
	proxyIdx    atomic.Int64
	proxyOnce   sync.Once
	proxyInited bool
)

func InitProxy() {
	proxyOnce.Do(func() {
		proxies := loadProxies()
		if len(proxies) == 0 {
			slog.Info("proxy_init", "status", "no_proxies_found")
			return
		}
		proxyList = proxies
		proxyIdx.Store(int64(rand.Intn(len(proxies))))
		http.DefaultTransport = &http.Transport{
			Proxy: rotateProxy,
		}
		proxyInited = true
		slog.Info("proxy_init", "count", len(proxies), "first", proxies[0].Host)
	})
}

func ProxyEnabled() bool { return proxyInited }

func rotateProxy(_ *http.Request) (*url.URL, error) {
	idx := proxyIdx.Add(1) - 1
	return proxyList[idx%int64(len(proxyList))], nil
}

func loadProxies() []*url.URL {
	if env := strings.TrimSpace(os.Getenv("CLOB_PROXY")); env != "" {
		switch strings.ToLower(env) {
		case "direct", "none", "off", "disabled":
			return nil
		}
		if u, err := url.Parse(env); err == nil {
			return []*url.URL{u}
		}
	}

	if env := os.Getenv("HTTPS_PROXY"); env != "" {
		if u, err := url.Parse(env); err == nil {
			return []*url.URL{u}
		}
	}
	return nil
}
