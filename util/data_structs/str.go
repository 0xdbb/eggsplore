package util

import (
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"strings"
	"unicode"
)

func IsEmptyStr(str string) bool {
	return len(str) == 0
}

func IsNotEmptyStr(str string) bool {
	return len(str) > 0
}

// WordToTitle - to title case
func WordToTitle(s string) string {
	return cases.Title(language.English).String(strings.ToLower(s))
}

// FilterEmptyStrings - returns a new slice of string without empty values
func FilterEmptyStrings(values []string) []string {
	var res = make([]string, 0, len(values))
	for _, v := range values {
		if IsNotEmptyStr(v) {
			res = append(res, v)
		}
	}
	return res
}

func TrimStrings(values []string) []string {
	for i, v := range values {
		values[i] = strings.TrimSpace(v)
	}
	return values
}

func Trim(text string) string {
	return strings.TrimSpace(text)
}

func TrimAll(items []string) []string {
	for i := range items {
		items[i] = Trim(items[i])
	}
	return items
}

func Tokenize(text string, separator ...string) []string {
	var sep = ","
	if len(separator) > 0 {
		sep = separator[0]
	}
	var tokens = strings.Split(text, sep)
	var results = make([]string, 0, len(tokens))
	for _, token := range TrimAll(strings.Split(text, ",")) {
		if len(token) > 0 {
			results = append(results, token)
		}
	}
	return results
}

func TokenizeMIDs(text string) []string {
	var tokens = strings.Split(text, ",")
	var results = make([]string, 0, len(tokens))
	for _, token := range TrimAll(strings.Split(text, ",")) {
		var subs = TrimAll(strings.Split(token, ":"))
		var tok string
		if len(subs) > 1 {
			tok = subs[1]
		} else if len(subs) == 1 {
			tok = subs[0]
		}
		if len(tok) > 30 {
			results = append(results, tok)
		}
	}
	return results
}

func AppendToken(src, token string) string {
	var tokens = Tokenize(src)
	tokens = append(tokens, token)
	return strings.Join(tokens, ",")
}

func AsCSV[T any](values []T, separator ...string) string {
	var sep = ","
	if len(separator) > 0 {
		sep = separator[0]
	}
	var items = make([]string, len(values))
	for i, v := range values {
		items[i] = fmt.Sprintf(`%v`, v)
	}
	return strings.Join(items, sep)
}

func CapitalizeStrings(arr []string, specials ...string) []string {
	capitalized := make([]string, len(arr))
	for i, str := range arr {
		capitalized[i] = Capitalize(str, specials...)
	}
	return capitalized
}

func Capitalize(str string, specialChars ...string) string {
	var capNext = true
	var specials = " "
	if len(specialChars) > 0 {
		specials = specialChars[0]
	}
	return strings.Map(func(r rune) rune {
		var char = string(r)
		if capNext && unicode.IsLetter(r) {
			r = unicode.ToUpper(r)
			capNext = false
		} else if strings.Contains(specials, char) {
			capNext = true
		} else {
			capNext = false
		}
		return r
	}, str)
}
