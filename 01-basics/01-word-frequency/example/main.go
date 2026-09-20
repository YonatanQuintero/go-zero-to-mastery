package main

import (
	"bufio"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {

	telemetryStream := "ORBITA: OK [Alt: 420km] 📡 Temp: -15°C 🛰️ STATUS: Estación_Alpha"

	scanner := bufio.NewScanner(strings.NewReader(telemetryStream))

	scanner.Split(bufio.ScanWords)

	fmt.Println("=== Procesando Tokens de Telemetría ===")

	for scanner.Scan() {
		word := scanner.Text()
		runeCount := utf8.RuneCountInString(word)
		byteCount := len(word)

		var specialSymbols []rune
		for _, r := range word {
			if !unicode.IsLetter(r) && !unicode.IsNumber(r) && !unicode.IsPunct(r) {
				specialSymbols = append(specialSymbols, r)
			}
		}

		fmt.Printf("Token: %-15s | Bytes: %2d | Runes %2d | Símbolos Especiales: %c\n",
			word, byteCount, runeCount, specialSymbols)
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error leyendo stream %v\n", err)
	}
}
