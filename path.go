package vrest

import "regexp"

var keysRegExp *regexp.Regexp = regexp.MustCompile(`\{[A-Za-z_0-9]+\}`)

func makePath(pathWithPlaceholders string, keysAndValues ...string) string {
	lenKeysAndValues := len(keysAndValues)
	if lenKeysAndValues == 0 || lenKeysAndValues%2 != 0 {
		return pathWithPlaceholders
	}

	return keysRegExp.ReplaceAllStringFunc(pathWithPlaceholders, func(s string) string {
		if len(s) < 3 {
			return s
		}
		key := s[1 : len(s)-1]
		for i := 0; i < lenKeysAndValues; i += 2 {
			if keysAndValues[i] == key {
				return keysAndValues[i+1]
			}
		}
		return s
	})
}
