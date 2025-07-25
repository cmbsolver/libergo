package lgstructs

import (
	runelib "characterrepo"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// DictionaryWord represents a structured word entry with various properties for linguistic and numerical analysis.
type DictionaryWord struct {
	DictionaryWordText          string `json:"dict_word"`
	RuneglishWordText           string `json:"dict_runeglish"`
	RuneWordText                string `json:"dict_rune"`
	RuneWordTextNoDoublet       string `json:"dict_rune_no_doublet"`
	GemSum                      int64  `json:"gem_sum"`
	GemSumPrime                 bool   `json:"gem_sum_prime"`
	GemProduct                  string `json:"gem_product"`
	GemProductPrime             bool   `json:"gem_product_prime"`
	DictionaryWordLength        int    `json:"dict_word_length"`
	RuneglishWordLength         int    `json:"dict_runeglish_length"`
	RuneWordLength              int    `json:"dict_rune_length"`
	DictRuneNoDoubletLength     int    `json:"dict_rune_no_doublet_length"`
	RunePattern                 string `json:"rune_pattern"`
	RunePatternNoDoubletPattern string `json:"rune_pattern_no_doublet"`
	RuneDistancePattern         string `json:"rune_distance_pattern"`
	Language                    string `json:"language"`
}

// GetRunePattern gets the rune pattern for the dictionary word
func GetRunePattern(word string) string {
	patternDictionary := make(map[int]string)
	var runes []string
	counter := 1

	for _, character := range word {
		if character == '\'' {
			runes = append(runes, "'")
			continue
		}

		found := false
		for key, value := range patternDictionary {
			if value == string(character) {
				runes = append(runes, fmt.Sprintf("%d", key))
				found = true
				break
			}
		}

		if !found {
			runes = append(runes, fmt.Sprintf("%d", counter))
			patternDictionary[counter] = string(character)
			counter++
		}
	}

	return strings.Join(runes, ",")
}

// RemoveDoublets removes consecutive duplicate characters from a word
func RemoveDoublets(word []string) string {
	if len(word) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString(word[0])

	for i := 1; i < len(word); i++ {
		if word[i] != word[i-1] {
			result.WriteString(word[i])
		}
	}

	return result.String()
}

// CalculateWordDistances calculates the distances between corresponding runes of two words based on their positions in runeList.
// It returns a comma-separated string of absolute distance values for each rune comparison.
func CalculateWordDistances(word1, word2, runeList []string) string {
	var result strings.Builder
	var isFirstTime = true

	for i := 0; i < len(word1); i++ {
		pos1 := getRuneIndex(word1[i], runeList)
		pos2 := getRuneIndex(word2[i], runeList)

		currentDistance := pos1 - pos2
		distance := int(math.Abs(float64(currentDistance)))

		if isFirstTime {
			result.WriteString(fmt.Sprintf("%s", strconv.Itoa(distance)))
		} else {
			result.WriteString(fmt.Sprintf(", %s", strconv.Itoa(distance)))
		}

		isFirstTime = false
	}

	return result.String()
}

// GetRuneDistancePattern generates a string representing the distances between adjacent runes in a slice of strings.
func GetRuneDistancePattern(word []string) string {
	charRepo := runelib.NewCharacterRepo()
	gemRunes := charRepo.GetGematriaRunes()

	if len(word) == 0 {
		return ""
	}

	var result strings.Builder
	currentValue := getRuneIndex(word[0], gemRunes)

	result.WriteString(strconv.Itoa(0))

	for i := 1; i < len(word); i++ {
		currentDistance := currentValue - getRuneIndex(word[i], gemRunes)
		distance := int(math.Abs(float64(currentDistance)))
		result.WriteString(fmt.Sprintf(", %s", strconv.Itoa(distance)))
		currentValue = getRuneIndex(word[i], gemRunes)
	}

	return result.String()
}

// GetRuneComparisonDistancePattern calculates a distance pattern between corresponding runes in two word slices.
func GetRuneComparisonDistancePattern(wordOne, wordTwo []string) string {
	charRepo := runelib.NewCharacterRepo()
	gemRunes := charRepo.GetGematriaRunes()

	if len(wordOne) == 0 {
		return ""
	}

	var result strings.Builder
	result.WriteString(strconv.Itoa(0))

	for i := 1; i < len(wordOne); i++ {
		currentDistance := getRuneIndex(wordOne[i], gemRunes) - getRuneIndex(wordTwo[i], gemRunes)
		distance := int(math.Abs(float64(currentDistance)))
		result.WriteString(fmt.Sprintf(", %s", strconv.Itoa(distance)))
	}

	return result.String()
}

// getRuneIndex finds the index of a given rune in the provided alphabet slice.
// Returns the index if found, otherwise returns -1.
func getRuneIndex(rune string, alphabet []string) int {
	for i, r := range alphabet {
		if r == rune {
			return i
		}
	}

	return -1
}
