// Package config holds tiny startup helpers that don't deserve their own pkg.
package config

import (
	"bufio"
	"os"
	"strings"
)

// LoadDotEnv reads simple KEY=VALUE lines from path and calls os.Setenv for
// any key not already set. Quiet on missing file (returns nil). Lines starting
// with '#' and empty lines are skipped. Values may be single- or double-quoted
// (quotes stripped). No shell expansion — keep it dumb.
//
// We use this so secrets in .env.local (TELEGRAM_BOT_TOKEN, later the
// wallet RPC, etc.) can be picked up by the bot without making scripts/
// alert-dispatch.sh the only consumer.
func LoadDotEnv(path string) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		eq := strings.IndexByte(line, '=')
		if eq <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		if len(val) >= 2 {
			first, last := val[0], val[len(val)-1]
			if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
				val = val[1 : len(val)-1]
			}
		}
		if _, already := os.LookupEnv(key); already {
			continue
		}
		_ = os.Setenv(key, val)
	}
	return sc.Err()
}
