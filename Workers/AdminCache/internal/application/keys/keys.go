package keys

import (
	"fmt"
	"strings"
	"unicode"
)

func normalizeRune(r rune) rune {
	switch r {
	case '\u00e1', '\u00e0', '\u00e4', '\u00e2', '\u00e3', '\u00e5':
		return 'a'
	case '\u00e9', '\u00e8', '\u00eb', '\u00ea':
		return 'e'
	case '\u00ed', '\u00ec', '\u00ef', '\u00ee':
		return 'i'
	case '\u00f3', '\u00f2', '\u00f6', '\u00f4', '\u00f5':
		return 'o'
	case '\u00fa', '\u00f9', '\u00fc', '\u00fb':
		return 'u'
	case '\u00f1':
		return 'n'
	case '\u00e7':
		return 'c'
	case '\u00df':
		return 's'
	default:
		return r
	}
}

func SlugCity(city string) string {
	s := strings.TrimSpace(strings.ToLower(city))
	if s == "" {
		return s
	}

	var b strings.Builder
	b.Grow(len(s))

	for _, ch := range s {
		ch = normalizeRune(ch)
		switch {
		case ch >= 'a' && ch <= 'z':
			b.WriteRune(ch)
		case ch >= '0' && ch <= '9':
			b.WriteRune(ch)
		case ch == '-' || ch == '_':
			b.WriteRune(ch)
		case ch == ' ':
			b.WriteRune('-')
		case unicode.IsSpace(ch):
			b.WriteRune('-')
		}
	}

	return b.String()
}

func RankGlobal(version string) string {
	return fmt.Sprintf("rank:global:%s", version)
}

func RankCity(citySlug, version string) string {
	return fmt.Sprintf("rank:city:%s:%s", citySlug, version)
}

func RankLockGlobal(version string) string {
	return fmt.Sprintf("rank:lock:global:%s", version)
}

func RankLockCity(citySlug, version string) string {
	return fmt.Sprintf("rank:lock:city:%s:%s", citySlug, version)
}

func CityIndex(version string) string {
	return fmt.Sprintf("rank:index:cities:%s", version)
}
