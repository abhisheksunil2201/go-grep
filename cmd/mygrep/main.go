package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "-E" {
		fmt.Fprintf(os.Stderr, "usage: echo <text> | go run cmd/mygrep/main.go -E \"<pattern>\"\n")
		os.Exit(2)
	}

	pattern := os.Args[2]

	line, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: read input text: %v\n", err)
		os.Exit(2)
	}

	ok, err := matchLine(string(line), pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(2)
	}

	if !ok {
		os.Exit(1)
	}
}

func matchLine(line string, pattern string) (bool, error) {
	fmt.Println("Logs from the program will appear here!")

	if pattern[0] == '^' {
		return match(line, pattern[1:], 0), nil
	}

	for i := 0; i <= len(line); i++ {
		if match(line, pattern, i) {
			return true, nil
		}
	}
	return false, nil
}

func match(line string, pattern string, curr int) bool {
	patternLength, stringLength, currPos := len(pattern), len(line), curr
	for i := 0; i < patternLength; i++ {
		if currPos >= stringLength {
			return pattern[i] == '$'
		}
		if i+1 < patternLength && pattern[i+1] == '?' {
			if line[currPos] == pattern[i] {
				if match(line, pattern[i+2:], currPos+1) {
					return true
				}
			}
			i++
			continue
		} else if pattern[i] == '\\' && i+1 < patternLength {
			if pattern[i+1] == 'd' && !unicode.IsDigit(rune(line[currPos])) {
				return false
			} else if pattern[i+1] == 'w' && !(unicode.IsLetter(rune(line[currPos])) || unicode.IsDigit(rune(line[currPos])) || line[currPos] == '_') {
				return false
			} else {
				i++
			}
		} else if pattern[i] == '[' && i+1 < patternLength {
			patternEndPos := strings.Index(pattern[i:], "]")
			patternHere := pattern[i+1 : patternEndPos]
			if pattern[i+1] == '^' {
				if strings.Contains(patternHere, string(line[currPos])) {
					return false
				}
			} else {
				if !strings.Contains(patternHere, string(line[currPos])) {
					return false
				}
			}
			i = patternEndPos
		} else if pattern[i] == '+' {
			if i == 0 {
				return false
			}
			for line[currPos] == pattern[i-1] && currPos < stringLength {
				currPos++
			}
			currPos-- //reset for next
		} else if pattern[i] == '.' {
			currPos++
			continue
		} else if pattern[i] == '(' {
			patternEndPos := strings.Index(pattern[i:], ")")
			innerPattern := pattern[i+1 : patternEndPos]
			if strings.Contains(innerPattern, "|") {
				patterns := strings.Split(innerPattern, "|")
				for _, p := range patterns {
					if match(line, p, currPos) {
						return true
					}
				}
				return false
			}
		} else {
			if line[currPos] != pattern[i] {
				return false
			}
		}
		currPos++
	}
	return true
}
