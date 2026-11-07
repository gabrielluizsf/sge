// Package i18n translates the words the example games show, in the
// language the player chose. English is the default, and the English text
// is the key of its translation: T("Save") is "Salvar" in Portuguese.
package i18n

import (
	"fmt"
	"sync/atomic"
)

// Language is a language the examples are shown in.
type Language int

// The languages, English first: it is the default.
const (
	English Language = iota
	Portuguese
)

var current atomic.Int32

// Set puts the language the games are shown in.
func Set(l Language) { current.Store(int32(l)) }

// Get is the language the games are shown in.
func Get() Language { return Language(current.Load()) }

// Next moves step languages along the list, and is the new one.
func Next(step int) Language {
	n := int(Portuguese) + 1
	Set(Language(((int(Get())+step)%n + n) % n))
	return Get()
}

// Name is the language's name, in the language itself.
func Name() string {
	if Get() == Portuguese {
		return "Português"
	}
	return "English"
}

// T is text in the language chosen, with args for a format string.
func T(text string, args ...any) string {
	return translated(text, text, args)
}

// TF is T for words said of a woman (woman) or of a man, for languages
// whose words change with it.
func TF(woman bool, text string, args ...any) string {
	key := text + "|m"
	if woman {
		key = text + "|f"
	}
	return translated(text, key, args)
}

// translated is English text, as it is, or its translation under key.
func translated(text, key string, args []any) string {
	if Get() == Portuguese {
		if s, ok := translations[key]; ok {
			text = s
		} else if s, ok := translations[text]; ok {
			text = s
		}
	}
	if len(args) > 0 {
		return fmt.Sprintf(text, args...)
	}
	return text
}
