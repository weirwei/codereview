package utils

import "testing"

func TestGetFileLanguageByFilepath(t *testing.T) {
	t.Run("TestGetFileLanguageByFilepath", func(t *testing.T) {
		filepath := "example.go"
		expectedLanguage := "go"
		actualLanguage := GetLangByFilepath(filepath)
		if actualLanguage != expectedLanguage {
			t.Errorf("Expected %s, got %s", expectedLanguage, actualLanguage)
		}
	})
}
