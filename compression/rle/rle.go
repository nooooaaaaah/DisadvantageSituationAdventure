package rle

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func compress(data string) string {
	var buffer bytes.Buffer
	prevChar := rune(data[0])
	count := 1

	for i := 1; i < len(data); i++ {
		if rune(data[i]) == prevChar {
			count++
		} else {
			buffer.WriteString(fmt.Sprintf("%c%d", prevChar, count))
			prevChar = rune(data[i])
			count = 1
		}
	}
	buffer.WriteString(fmt.Sprintf("%c%d", prevChar, count))

	return buffer.String()
}

// Note: When converting a rune/byte to int convert to string first
func decompress(data string) string {
	var buffer bytes.Buffer
	prevChar := rune(data[0])

	for i := 1; i < len(data); i++ {
		if unicode.IsDigit(rune(data[i])) {
			count, err := strconv.Atoi(string(data[i]))
			if err != nil {
				fmt.Println("Uh oh looks like you suck: ", err)
			}
			buffer.WriteString(strings.Repeat(string(prevChar), count))
		} else {
			prevChar = rune(data[i])
		}
	}

	return buffer.String()
}

func main() {
	input := "aaaabbbccddddd"
	compressed := compress(input)
	decompressed := decompress(compressed)
	fmt.Println("Original:", input)
	fmt.Println("Compressed:", compressed)
	fmt.Println("Decompressed:", decompressed)
}
