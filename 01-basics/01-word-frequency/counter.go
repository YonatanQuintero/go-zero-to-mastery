package wordfrequency

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

func WordFrequency(r io.Reader) (map[string]int, error) {

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanWords)
	counts := make(map[string]int)

	for scanner.Scan() {
		t := strings.ToLower(scanner.Text())
		w := strings.TrimFunc(t, func(r rune) bool {
			return !unicode.IsPunct(r)
		})

		if len(strings.TrimSpace(w)) > 0 {
			counts[w]++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading stream %v\n", err)
		return nil, err
	}

	return counts, nil
}

func CharFrequency(r io.Reader) (map[rune]int, error) {

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanRunes)
	counts := make(map[rune]int)

	for scanner.Scan() {
		t := strings.ToLower(scanner.Text())
		c := []rune(t)[0]
		if !unicode.IsPunct(c) && !unicode.IsSpace(c) {
			counts[c]++
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading stream %v\n", err)
		return nil, err
	}

	return counts, nil
}
