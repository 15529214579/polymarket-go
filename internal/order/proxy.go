package order

import (
	"bufio"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

var (
	proxyList  []*url.URL
	proxyIdx   atomic.Int64
	proxyOnce  sync.Once
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
	if env := os.Getenv("CLOB_PROXY"); env != "" {
		if u, err := url.Parse(env); err == nil {
			return []*url.URL{u}
		}
	}

	if env := os.Getenv("HTTPS_PROXY"); env != "" {
		if u, err := url.Parse(env); err == nil {
			return []*url.URL{u}
		}
	}

	paths := []string{"proxies.txt"}
	if exe, err := os.Executable(); err == nil {
		paths = append(paths, filepath.Join(filepath.Dir(exe), "..", "proxies.txt"))
	}
	for _, p := range paths {
		if list := readProxyFile(p); len(list) > 0 {
			return list
		}
	}
	return nil
}

func readProxyFile(path string) []*url.URL {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	var out []*url.URL
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		u, err := url.Parse(line)
		if err != nil {
			continue
		}
		out = append(out, u)
	}
	return out
}
