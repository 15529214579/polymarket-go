// Package sanitize redacts credentials from errors and logs.
package sanitize

import "regexp"

var (
	telegramBotTokenRE = regexp.MustCompile(`/bot[0-9]{8,10}:[A-Za-z0-9_-]{20,}`)
	apiKeyQueryRE      = regexp.MustCompile(`([?&]apiKey=)[^&"\s]+`)
)

// SecretString redacts common credential shapes embedded in error strings.
func SecretString(s string) string {
	s = telegramBotTokenRE.ReplaceAllString(s, `/bot<redacted>`)
	s = apiKeyQueryRE.ReplaceAllString(s, `${1}<redacted>`)
	return s
}

func Error(err error) string {
	if err == nil {
		return ""
	}
	return SecretString(err.Error())
}
