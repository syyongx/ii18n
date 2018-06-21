package ii18n

import (
	"strings"
	"errors"
)

var patternNumeric = "(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"

type Formatter struct {
}

func NewFormatter() *Formatter {
	return &Formatter{}
}

func (f *Formatter) format(pattern string, params map[string]string, lang string) (string, error) {
	tokens := f.tokenizePattern(pattern)
	if tokens == nil {
		return "", errors.New("message pattern is invalid")
	}

	return strings.Join(tokens, ""), nil
}

// Tokenizes a pattern by separating normal text from replaceable patterns.
func (f *Formatter) tokenizePattern(pattern string) []string {
	pos := strings.Index(pattern, "{")
	if pos == -1 {
		return []string{pattern}
	}
	//pr := []rune(pattern)
	start, depth, length := pos, 1, len(pattern)
	tokens := []string{pattern[:pos]}
	for {
		if pos+1 > length {
			break
		}
		open := strings.Index(pattern[pos+1:], "{")
		close := strings.Index(pattern[pos+1:], "}")
		if open == -1 && close == -1 {
			break
		}
		if open == -1 {
			open = length
		}
		if close > open {
			depth++
			pos = open
		} else {
			depth--
			pos = close
		}
		if depth == 0 {
			start = pos + 1
			tokens = append(tokens, pattern[start:open])
			start = open
		}
		if depth != 0 && (open == -1 || close == -1) {
			break
		}
	}
	if depth != 0 {
		return nil
	}

	return tokens
}
