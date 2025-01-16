package main

import (
	"fmt"

	"github.com/YakDriver/master-sifo-dyas/sifo"
)

func main() {
	words := sifo.LoadWords("words_with_rank.csv")

	if 1 == 0 {
		partialWords := sifo.CreatePartialWordDictionary(words)
		if err := sifo.WritePartialWordsToCSV(partialWords, "partial_words.csv"); err != nil {
			fmt.Printf("Error writing to CSV: %v\n", err)
		} else {
			fmt.Println("Partial word dictionary written to partial_words.csv")
		}
	}

	bestCipher := sifo.FindBestCipher(words, 20000)
	sifo.CountHighScore(words, bestCipher, true)

	fmt.Println("Best Cipher:")
	for k, v := range bestCipher {
		fmt.Printf("%s -> %s\n", k, string(v))
	}
}
