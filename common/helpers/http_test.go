package helpers

import (
	"net/http/httptest"
	"testing"
)

func TestParseAcceptLanguage_Empty(t *testing.T) {
	got := ParseAcceptLanguage("", []string{"en", "sl", "de"})
	if got != "en" {
		t.Errorf("empty header: expected %q, got %q", "en", got)
	}
}

func TestParseAcceptLanguage_SingleMatch(t *testing.T) {
	got := ParseAcceptLanguage("de", []string{"en", "sl", "de"})
	if got != "de" {
		t.Errorf("single match: expected %q, got %q", "de", got)
	}
}

func TestParseAcceptLanguage_MultipleValues(t *testing.T) {
	got := ParseAcceptLanguage("fr, de, en", []string{"en", "sl", "de"})
	if got != "de" {
		t.Errorf("multiple values: expected %q, got %q", "de", got)
	}
}

func TestParseAcceptLanguage_WithQuality(t *testing.T) {
	got := ParseAcceptLanguage("sl-SI;q=0.9, en;q=0.8", []string{"en", "sl", "de"})
	if got != "sl" {
		t.Errorf("with quality: expected %q, got %q", "sl", got)
	}
}

func TestParseAcceptLanguage_NoMatch(t *testing.T) {
	got := ParseAcceptLanguage("fr, ja", []string{"en", "sl", "de"})
	if got != "en" {
		t.Errorf("no match should default to en: expected %q, got %q", "en", got)
	}
}

func TestParseAcceptLanguage_RegionCode(t *testing.T) {
	got := ParseAcceptLanguage("en-US, sl", []string{"en", "sl", "de"})
	if got != "en" {
		t.Errorf("region code: expected %q, got %q", "en", got)
	}
}

func TestParseAcceptLanguage_LowercaseNormalization(t *testing.T) {
	got := ParseAcceptLanguage("DE", []string{"en", "sl", "de"})
	if got != "de" {
		t.Errorf("lowercase normalization: expected %q, got %q", "de", got)
	}
}

func TestParseAcceptLanguage_FirstMatchWins(t *testing.T) {
	got := ParseAcceptLanguage("de, sl, en", []string{"en", "sl", "de"})
	if got != "de" {
		t.Errorf("first match should win: expected %q, got %q", "de", got)
	}
}

func TestParsePagination_Defaults(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/items", nil)
	pn, ps := ParsePagination(req, 50)

	if pn != 1 {
		t.Errorf("default page number: expected 1, got %d", pn)
	}
	if ps != 50 {
		t.Errorf("default page size: expected 50, got %d", ps)
	}
}

func TestParsePagination_CustomValues(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/items?page[number]=3&page[size]=25", nil)
	pn, ps := ParsePagination(req, 50)

	if pn != 3 {
		t.Errorf("custom page number: expected 3, got %d", pn)
	}
	if ps != 25 {
		t.Errorf("custom page size: expected 25, got %d", ps)
	}
}

func TestParsePagination_InvalidValues(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/items?page[number]=abc&page[size]=-5", nil)
	pn, ps := ParsePagination(req, 50)

	if pn != 1 {
		t.Errorf("invalid page number falls back to 1: expected 1, got %d", pn)
	}
	if ps != 50 {
		t.Errorf("invalid page size falls back to default: expected 50, got %d", ps)
	}
}

func TestParsePagination_ZeroValues(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/items?page[number]=0&page[size]=0", nil)
	pn, ps := ParsePagination(req, 50)

	if pn != 1 {
		t.Errorf("zero page number falls back to 1: expected 1, got %d", pn)
	}
	if ps != 50 {
		t.Errorf("zero page size falls back to default: expected 50, got %d", ps)
	}
}

func TestParsePagination_NegativeValues(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/items?page[number]=-1&page[size]=-10", nil)
	pn, ps := ParsePagination(req, 50)

	if pn != 1 {
		t.Errorf("negative page number falls back to 1: expected 1, got %d", pn)
	}
	if ps != 50 {
		t.Errorf("negative page size falls back to default: expected 50, got %d", ps)
	}
}

func TestParsePagination_CustomDefault(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/items", nil)
	pn, ps := ParsePagination(req, 20)

	if pn != 1 {
		t.Errorf("custom default page number: expected 1, got %d", pn)
	}
	if ps != 20 {
		t.Errorf("custom default page size: expected 20, got %d", ps)
	}
}

func TestParsePagination_OnlyPageNumber(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/items?page[number]=5", nil)
	pn, ps := ParsePagination(req, 50)

	if pn != 5 {
		t.Errorf("page number: expected 5, got %d", pn)
	}
	if ps != 50 {
		t.Errorf("page size should use default: expected 50, got %d", ps)
	}
}
