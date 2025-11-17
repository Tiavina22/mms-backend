package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// I18n holds translation data for multiple languages
type I18n struct {
	translations map[string]map[string]string
	mu           sync.RWMutex
	defaultLang  string
}

var (
	i18nInstance *I18n
	once         sync.Once
)

// GetI18n returns the singleton instance of I18n
func GetI18n() *I18n {
	once.Do(func() {
		i18nInstance = &I18n{
			translations: make(map[string]map[string]string),
			defaultLang:  "en",
		}
		i18nInstance.LoadTranslations()
	})
	return i18nInstance
}

// LoadTranslations loads all translation files from the locales directory
func (i *I18n) LoadTranslations() {
	i.mu.Lock()
	defer i.mu.Unlock()

	localesDir := "locales"
	
	// Check if locales directory exists
	if _, err := os.Stat(localesDir); os.IsNotExist(err) {
		log.Printf("Warning: locales directory not found at %s", localesDir)
		return
	}

	files, err := os.ReadDir(localesDir)
	if err != nil {
		log.Printf("Error reading locales directory: %v", err)
		return
	}

	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// Extract language code from filename (e.g., "en.json" -> "en")
		lang := strings.TrimSuffix(file.Name(), ".json")
		
		filePath := filepath.Join(localesDir, file.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			log.Printf("Error reading translation file %s: %v", file.Name(), err)
			continue
		}

		var translations map[string]string
		if err := json.Unmarshal(data, &translations); err != nil {
			log.Printf("Error parsing translation file %s: %v", file.Name(), err)
			continue
		}

		i.translations[lang] = translations
		log.Printf("Loaded translations for language: %s", lang)
	}
}

// Translate returns the translated message for the given key and language
func (i *I18n) Translate(lang, key string, args ...interface{}) string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	// Try to get translation for the specified language
	if langMap, exists := i.translations[lang]; exists {
		if translation, exists := langMap[key]; exists {
			// If args are provided, format the translation
			if len(args) > 0 {
				return fmt.Sprintf(translation, args...)
			}
			return translation
		}
	}

	// Fallback to default language
	if lang != i.defaultLang {
		if langMap, exists := i.translations[i.defaultLang]; exists {
			if translation, exists := langMap[key]; exists {
				if len(args) > 0 {
					return fmt.Sprintf(translation, args...)
				}
				return translation
			}
		}
	}

	// If no translation found, return the key itself
	return key
}

// T is a shorthand for Translate
func T(lang, key string, args ...interface{}) string {
	return GetI18n().Translate(lang, key, args...)
}

// GetLanguageFromPhone attempts to determine language from phone prefix
func GetLanguageFromPhone(phone string) string {
	// Remove any spaces or special characters
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "(", "")
	phone = strings.ReplaceAll(phone, ")", "")

	// Map of country codes to languages
	prefixToLang := map[string]string{
		"+1":   "en", // US/Canada
		"+33":  "fr", // France
		"+44":  "en", // UK
		"+49":  "de", // Germany
		"+34":  "es", // Spain
		"+39":  "it", // Italy
		"+86":  "zh", // China
		"+81":  "ja", // Japan
		"+82":  "ko", // South Korea
		"+91":  "hi", // India
		"+55":  "pt", // Brazil
		"+7":   "ru", // Russia
		"+261": "fr", // Madagascar (French-speaking)
	}

	// Check for matching prefix
	for prefix, lang := range prefixToLang {
		if strings.HasPrefix(phone, prefix) {
			return lang
		}
	}

	// Default to English
	return "en"
}

// SupportedLanguages returns a list of supported language codes
func (i *I18n) SupportedLanguages() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	langs := make([]string, 0, len(i.translations))
	for lang := range i.translations {
		langs = append(langs, lang)
	}
	return langs
}
