package keys

import (
	"fmt"
	"strings"
)

func SlugCity(city string) string {
	city = strings.ToLower(city)
	city = strings.ReplaceAll(city, " ", "-")
	city = strings.ReplaceAll(city, "á", "a")
	city = strings.ReplaceAll(city, "é", "e")
	city = strings.ReplaceAll(city, "í", "i")
	city = strings.ReplaceAll(city, "ó", "o")
	city = strings.ReplaceAll(city, "ú", "u")
	city = strings.ReplaceAll(city, "ä", "a")
	city = strings.ReplaceAll(city, "ë", "e")
	city = strings.ReplaceAll(city, "ï", "i")
	city = strings.ReplaceAll(city, "ö", "o")
	city = strings.ReplaceAll(city, "ü", "u")
	return city
}

func RankGlobal(page, size int) string {
	return fmt.Sprintf("rank:global:page:%d:size:%d", page, size)
}
func RankCity(citySlug string, page, size int) string {
	return fmt.Sprintf("rank:city:%s:page:%d:size:%d", citySlug, page, size)
}
