package keys

import (
	"fmt"
	"strings"
)

func SlugCity(city string) string {
	s := strings.TrimSpace(strings.ToLower(city))
	if s == "" {
		return s
	}
	// Reemplazos comunes de tildes y ñ/Ñ
	replacer := strings.NewReplacer(
		"á", "a", "é", "e", "í", "i", "ó", "o", "ú", "u",
		"à", "a", "è", "e", "ì", "i", "ò", "o", "ù", "u",
		"ä", "a", "ë", "e", "ï", "i", "ö", "o", "ü", "u",
		"â", "a", "ê", "e", "î", "i", "ô", "o", "û", "u",
		"Á", "a", "É", "e", "Í", "i", "Ó", "o", "Ú", "u",
		"Ä", "a", "Ë", "e", "Ï", "i", "Ö", "o", "Ü", "u",
		"Â", "a", "Ê", "e", "Î", "i", "Ô", "o", "Û", "u",
		"ñ", "n", "Ñ", "n",
	)
	s = replacer.Replace(s)
	// Espacios a '-'
	s = strings.ReplaceAll(s, " ", "-")
	// Mantener solo [a-z0-9_-]
	var b strings.Builder
	for _, ch := range s {
		if (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '-' || ch == '_' {
			b.WriteRune(ch)
		}
	}
	return b.String()
}

func RankGlobal(page, size int) string {
	return fmt.Sprintf("rank:global:page:%d:size:%d", page, size)
}
func RankCity(citySlug string, page, size int) string {
	return fmt.Sprintf("rank:city:%s:page:%d:size:%d", citySlug, page, size)
}

