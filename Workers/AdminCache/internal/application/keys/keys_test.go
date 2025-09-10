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
        {name: "city with accents", input: "São Paulo", expected: "sao-paulo"},
        {name: "city with spanish accents", input: "Córdoba", expected: "cordoba"},
        {name: "city with multiple accents", input: "Bogotá Medellín", expected: "bogota-medellin"},
        {name: "city with german umlauts", input: "München", expected: "munchen"},
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
    tests := []struct {
        name     string
        page     int
        size     int
        expected string
    }{
        {name: "first page", page: 1, size: 10, expected: "rank:global:page:1:size:10"},
        {name: "second page", page: 2, size: 20, expected: "rank:global:page:2:size:20"},
        {name: "zero values", page: 0, size: 0, expected: "rank:global:page:0:size:0"},
        {name: "large values", page: 100, size: 50, expected: "rank:global:page:100:size:50"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := RankGlobal(tt.page, tt.size)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestRankCity(t *testing.T) {
    tests := []struct {
        name     string
        citySlug string
        page     int
        size     int
        expected string
    }{
        {name: "madrid first page", citySlug: "madrid", page: 1, size: 10, expected: "rank:city:madrid:page:1:size:10"},
        {name: "new-york second page", citySlug: "new-york", page: 2, size: 20, expected: "rank:city:new-york:page:2:size:20"},
        {name: "empty city slug", citySlug: "", page: 1, size: 10, expected: "rank:city::page:1:size:10"},
        {name: "complex city slug", citySlug: "sao-paulo", page: 5, size: 25, expected: "rank:city:sao-paulo:page:5:size:25"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := RankCity(tt.citySlug, tt.page, tt.size)
            assert.Equal(t, tt.expected, result)
        })
    }
}

func TestSlugCity_Integration(t *testing.T) {
    // Test integration between SlugCity and RankCity
    cityName := "São Paulo"
    slug := SlugCity(cityName)
    rankKey := RankCity(slug, 1, 10)

    expected := "rank:city:sao-paulo:page:1:size:10"
    assert.Equal(t, expected, rankKey)
}

