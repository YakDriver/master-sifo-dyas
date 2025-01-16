package sifo

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func LoadWords(filename string) map[string]int64 {
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	words := make(map[string]int64)
	scanner := bufio.NewScanner(file)
	i := 0
	for scanner.Scan() {
		i++
		if i == 1 {
			continue // skip first line
		}
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) < 3 {
			continue // skip lines that don't have at least 3 parts
		}
		value, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			continue // skip lines where the third part is not a valid int64
		}
		words[parts[1]] = value
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}

	fmt.Printf("Loaded %d words\n", len(words))
	return words
}

func CreatePartialWordDictionary(words map[string]int64) map[string]int64 {
	partialWords := make(map[string]int64)

	for word, count := range words {
		length := len(word)
		for i := 0; i < length; i++ {
			for j := i + 1; j <= length; j++ {
				substr := word[i:j]
				partialWords[substr] += count
			}
		}
	}

	return partialWords
}

func WritePartialWordsToCSV(partialWords map[string]int64, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"Partial Word", "Count"}); err != nil {
		return err
	}

	// Write data
	for k, v := range partialWords {
		if err := writer.Write([]string{k, strconv.FormatInt(v, 10)}); err != nil {
			return err
		}
	}

	return nil
}
