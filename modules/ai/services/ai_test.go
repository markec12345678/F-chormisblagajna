package services

import (
	"testing"
)

func TestFuzzyMatch_ExactMatch(t *testing.T) {
	score := FuzzyMatch("pizza", "pizza")
	if score != 1.0 {
		t.Errorf("expected 1.0 for exact match, got %f", score)
	}
}

func TestFuzzyMatch_ExactMatchCaseInsensitive(t *testing.T) {
	score := FuzzyMatch("Pizza", "pizza")
	if score != 1.0 {
		t.Errorf("expected 1.0 for case-insensitive match, got %f", score)
	}
}

func TestFuzzyMatch_Prefix(t *testing.T) {
	score := FuzzyMatch("piz", "pizza")
	if score < 0.85 || score > 0.95 {
		t.Errorf("expected ~0.9 for prefix match, got %f", score)
	}
}

func TestFuzzyMatch_Contains(t *testing.T) {
	score := FuzzyMatch("izz", "pizza")
	if score < 0.7 || score > 0.85 {
		t.Errorf("expected ~0.8 for contains match, got %f", score)
	}
}

func TestFuzzyMatch_Typo(t *testing.T) {
	score := FuzzyMatch("piza", "pizza")
	if score < 0.3 {
		t.Errorf("expected >0.3 for typo match, got %f", score)
	}
}

func TestFuzzyMatch_SlovenianChars(t *testing.T) {
	score := FuzzyMatch("čevapčiči", "cevapcici")
	if score != 1.0 {
		t.Errorf("expected 1.0 for Slovenian char normalization, got %f", score)
	}
}

func TestFuzzyMatch_SlovenianCharsPartial(t *testing.T) {
	score := FuzzyMatch("š", "šunka")
	if score < 0.7 {
		t.Errorf("expected >0.7 for Slovenian char prefix, got %f", score)
	}
}

func TestFuzzyMatch_SlovenianCharsContain(t *testing.T) {
	score := FuzzyMatch("ž", "žara")
	if score >= 0.7 && score <= 1.0 {
		t.Logf("Slovenian ž in žara: %f", score)
	} else {
		t.Errorf("unexpected score for Slovenian char: %f", score)
	}
}

func TestFuzzyMatch_MultiWord(t *testing.T) {
	score := FuzzyMatch("ice cream", "ice cream cake")
	if score >= 0.5 {
		t.Logf("multi-word score: %f", score)
	} else {
		t.Errorf("expected >=0.5 for multi-word partial, got %f", score)
	}
}

func TestFuzzyMatch_EmptyStrings(t *testing.T) {
	score := FuzzyMatch("", "")
	if score != 0 {
		t.Errorf("expected 0 for empty strings, got %f", score)
	}
}

func TestFuzzyMatch_EmptyQuery(t *testing.T) {
	score := FuzzyMatch("", "pizza")
	if score != 0 {
		t.Errorf("expected 0 for empty query, got %f", score)
	}
}

func TestFuzzyMatch_CompletelyDifferent(t *testing.T) {
	score := FuzzyMatch("xyz", "pizza")
	if score > 0.3 {
		t.Errorf("expected <0.3 for completely different words, got %f", score)
	}
}

func TestSoundex_SameForSimilar(t *testing.T) {
	s1 := soundex("peter")
	s2 := soundex("petar")
	if s1 != s2 {
		t.Errorf("expected same soundex for peter/petar, got %s and %s", s1, s2)
	}
}

func TestSoundex_Slovenian(t *testing.T) {
	s1 := soundex("čevapčiči")
	s2 := soundex("cevapcici")
	if s1 != s2 {
		t.Errorf("expected same soundex for Slovenian chars, got %s and %s", s1, s2)
	}
}

func TestSoundex_Empty(t *testing.T) {
	s := soundex("")
	if s != "" {
		t.Errorf("expected empty string for empty input, got %s", s)
	}
}

func TestLevenshteinDistance_Same(t *testing.T) {
	d := levenshteinDistance("pizza", "pizza")
	if d != 0 {
		t.Errorf("expected 0 for same strings, got %d", d)
	}
}

func TestLevenshteinDistance_OneEdit(t *testing.T) {
	d := levenshteinDistance("pizza", "piza")
	if d != 1 {
		t.Errorf("expected 1 for one edit, got %d", d)
	}
}

func TestLevenshteinDistance_Slovenian(t *testing.T) {
	d := levenshteinDistance("čevap", "cevap")
	if d != 0 {
		t.Errorf("expected 0 for Slovenian normalization, got %d", d)
	}
}

func TestLevenshteinDistance_Empty(t *testing.T) {
	d := levenshteinDistance("", "abc")
	if d != 3 {
		t.Errorf("expected 3 for empty vs abc, got %d", d)
	}
}

func TestNormalizeSlovenian(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"čevapčiči", "cevapcici"},
		{"Šunka", "Sunka"},
		{"Žar", "Zar"},
		{"piza", "piza"},
		{"", ""},
	}

	for _, tt := range tests {
		result := normalizeSlovenian(tt.input)
		if result != tt.expected {
			t.Errorf("normalizeSlovenian(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestFuzzyMatch_ScoreRanges(t *testing.T) {
	tests := []struct {
		name      string
		query     string
		target    string
		minScore  float64
		maxScore  float64
	}{
		{"exact", "coffee", "coffee", 1.0, 1.0},
		{"prefix", "cof", "coffee", 0.85, 0.95},
		{"contains", "off", "coffee", 0.7, 0.85},
		{"typo", "cofee", "coffee", 0.3, 0.7},
		{"similar_sound", "kofe", "coffee", 0.0, 0.6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := FuzzyMatch(tt.query, tt.target)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("FuzzyMatch(%q, %q) = %f, want [%f, %f]", tt.query, tt.target, score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestFuzzyMatch_SlovenianCompleteWords(t *testing.T) {
	score := FuzzyMatch("solata", "Solata")
	if score != 1.0 {
		t.Errorf("expected 1.0 for Slovenian case-insensitive match, got %f", score)
	}
}

func TestFuzzyMatch_SlovenianPartialWords(t *testing.T) {
	score := FuzzyMatch("solat", "solata")
	if score < 0.85 {
		t.Errorf("expected >0.85 for Slovenian prefix, got %f", score)
	}
}
