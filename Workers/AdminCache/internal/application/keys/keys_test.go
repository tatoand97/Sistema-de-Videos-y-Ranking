package keys

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlugCity(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "simple city name", input: "Madrid", expected: "madrid"},
		{name: "city with spaces", input: "New York", expected: "new-york"},
		{name: "city with accents", input: "S\u00e3o Paulo", expected: "sao-paulo"},
		{name: "city with spanish accents", input: "C\u00f3rdoba", expected: "cordoba"},
		{name: "city with multiple accents", input: "Bogot\u00e1 Medell\u00edn", expected: "bogota-medellin"},
		{name: "city with german umlauts", input: "M\u00fcnchen", expected: "munchen"},
		{name: "empty string", input: "", expected: ""},
		{name: "already lowercase", input: "barcelona", expected: "barcelona"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SlugCity(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestRankGlobal(t *testing.T) {
	assert.Equal(t, "rank:global:v2", RankGlobal("v2"))
	assert.Equal(t, "rank:global:v2025", RankGlobal("v2025"))
}

func TestRankCity(t *testing.T) {
	assert.Equal(t, "rank:city:madrid:v2", RankCity("madrid", "v2"))
	assert.Equal(t, "rank:city::v2", RankCity("", "v2"))
	assert.Equal(t, "rank:city:sao-paulo:v9", RankCity("sao-paulo", "v9"))
}

func TestLockKeys(t *testing.T) {
	assert.Equal(t, "rank:lock:global:v2", RankLockGlobal("v2"))
	assert.Equal(t, "rank:lock:city:madrid:v2", RankLockCity("madrid", "v2"))
}

func TestCityIndex(t *testing.T) {
	assert.Equal(t, "rank:index:cities:v2", CityIndex("v2"))
}

func TestSlugCityIntegration(t *testing.T) {
	cityName := "S\u00e3o Paulo"
	slug := SlugCity(cityName)
	rankKey := RankCity(slug, "v2")

	assert.Equal(t, "rank:city:sao-paulo:v2", rankKey)
}
