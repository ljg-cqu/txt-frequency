package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/sqweek/dialog"
)

/*
Description:
This program analyzes text files to extract and categorize Chinese characters, Chinese words, English words, and English phrases, providing both deduplicated and duplicated outputs.

Features:
- GUI-based file selection for ease of use.
- Categorizes text into Chinese characters, Chinese words, English words, and English phrases.
- Generates frequency-based deduplicated outputs and preserves original order for duplicated elements.
- Supports regex-based text processing and sorting by frequency.

Workflow:
1. Users select an input file via a GUI dialog.
2. The program reads the input, categorizing Chinese and English text using regex patterns:
   - Chinese characters and words.
   - English words and phrases.
3. Frequency maps and original lists are constructed for text elements.
4. Deduplicated text is sorted by frequency and saved to corresponding output files:
   - `deduplicated_chinese.txt` and `deduplicated_english.txt`.
5. Raw duplicated data is saved preserving original order:
   - `duplicated_chinese.txt` and `duplicated_english.txt`.
6. All outputs are written and saved with success notifications.
*/

func main() {
	// Allow users to specify the input file
	fmt.Println("Select the input file:")
	inputFile, err := dialog.File().
		Title("Select Input File").
		Filter("Text Files (*.txt)", "txt").
		Load()
	if err != nil {
		fmt.Printf("Error selecting input file: %v\n", err)
		return
	}
	if inputFile == "" {
		fmt.Println("No input file selected.")
		return
	}
	fmt.Printf("Selected input file: %s\n", inputFile)

	// Predefined output files
	chineseFileDedup := "deduplicated_chinese.txt"
	chineseFileDup := "duplicated_chinese.txt"
	englishFileDedup := "deduplicated_english.txt"
	englishFileDup := "duplicated_english.txt"

	// Open the input file
	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("Error opening input file: %v\n", err)
		return
	}
	defer file.Close()

	// Regex patterns
	chineseCharacterRegex := `[\p{Han}]`                         // Matches individual Chinese characters
	chineseWordsRegex := `[\p{Han}]+`                            // Matches sequences of Chinese characters as words
	englishWordRegex := `\b[a-zA-Z0-9']+(?:-[a-zA-Z0-9']+)?\b`   // Matches English words and compounds like "micro-video", also handle "I'll"
	englishPhrasesRegex := `\b[a-zA-Z0-9][\w\s'-]*[a-zA-Z0-9]\b` // Matches English phrases with spaces

	// Frequency maps
	chineseCharFreq := make(map[string]int)
	chineseWordsFreq := make(map[string]int)
	englishWordFreq := make(map[string]int)
	englishPhrasesFreq := make(map[string]int)

	// Lists to retain duplications (as they appear in the original order)
	chineseCharList := []string{}
	chineseWordsList := []string{}
	englishWordList := []string{}
	englishPhrasesList := []string{}

	// Read the input file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Match and process Chinese characters
		chineseCharMatches := regexp.MustCompile(chineseCharacterRegex).FindAllString(line, -1)
		for _, char := range chineseCharMatches {
			chineseCharFreq[char]++
			chineseCharList = append(chineseCharList, char) // Append in original order
		}

		// Match and process Chinese words
		chineseWordMatches := regexp.MustCompile(chineseWordsRegex).FindAllString(line, -1)
		for _, word := range chineseWordMatches {
			chineseWordsFreq[word]++
			chineseWordsList = append(chineseWordsList, word) // Append in original order
		}

		// Match and process English words (with hyphenated compounds like "micro-video")
		englishWordMatches := regexp.MustCompile(englishWordRegex).FindAllString(line, -1)
		for _, word := range englishWordMatches {
			normalizedWord := strings.ToLower(word) // Normalize to lowercase for consistency
			englishWordFreq[normalizedWord]++
			englishWordList = append(englishWordList, word) // Append in original order
		}

		// Match and process English phrases
		englishPhraseMatches := regexp.MustCompile(englishPhrasesRegex).FindAllString(line, -1)
		for _, phrase := range englishPhraseMatches {
			normalizedPhrase := strings.ToLower(strings.TrimSpace(phrase)) // Normalize case and trim
			englishPhrasesFreq[normalizedPhrase]++
			englishPhrasesList = append(englishPhrasesList, phrase) // Append in original order
		}
	}

	// Handle scanner error
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading input file: %v\n", err)
		return
	}

	// Sort lists by frequency (descending order) for deduplicated outputs
	chineseCharDedupSorted := sortByFrequency(chineseCharFreq)
	englishWordDedupSorted := sortByFrequency(englishWordFreq)

	// Write output files
	writeToFile(chineseFileDedup, chineseCharDedupSorted) // Deduplicated Chinese characters
	writeToFile(chineseFileDup, chineseCharList)          // Duplicated Chinese characters (original order)
	writeToFile(englishFileDedup, englishWordDedupSorted) // Deduplicated English words
	writeToFile(englishFileDup, englishWordList)          // Duplicated English words (original order)

	fmt.Println("All output files written successfully.")
}

// Function to write data to a file
func writeToFile(filePath string, data []string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, item := range data {
		writer.WriteString(item + "\n")
	}
	writer.Flush()
}

// Helper function to sort map entries by frequency (descending order)
func sortByFrequency(freqMap map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}

	// Create a slice of key-value pairs
	var sortedPairs []kv
	for k, v := range freqMap {
		sortedPairs = append(sortedPairs, kv{k, v})
	}

	// Sort by frequency in descending order
	sort.Slice(sortedPairs, func(i, j int) bool {
		return sortedPairs[i].Value > sortedPairs[j].Value
	})

	// Extract sorted keys
	var sortedKeys []string
	for _, pair := range sortedPairs {
		sortedKeys = append(sortedKeys, pair.Key)
	}

	return sortedKeys
}
