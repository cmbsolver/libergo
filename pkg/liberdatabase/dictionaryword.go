package liberdatabase

import (
	"fmt"
	"gorm.io/gorm"
	"runer"
	"strings"
)

type DictionaryWord struct {
	gorm.Model
	DictionaryWordText   string `gorm:"column:dict_word"`
	RuneglishWordText    string `gorm:"column:dict_runeglish"`
	RuneWordText         string `gorm:"column:dict_rune"`
	GemSum               int64  `gorm:"column:gem_sum"`
	DictionaryWordLength int    `gorm:"column:dict_word_length"`
	RuneglishWordLength  int    `gorm:"column:dict_runeglish_length"`
	RuneWordLength       int    `gorm:"column:dict_rune_length"`
	RunePattern          string `gorm:"column:rune_pattern"`
}

func (DictionaryWord) TableName() string {
	return "public.dictionary_words"
}

// GetRunePattern gets the rune pattern for the dictionary word
func (dw DictionaryWord) GetRunePattern() string {
	patternDictionary := make(map[int]string)
	var runes []string
	counter := 1

	for _, character := range dw.RuneWordText {
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

// String returns a string representation of the dictionary word
func (dw DictionaryWord) String() string {
	return fmt.Sprintf("%s - %s - %s - %d", dw.DictionaryWordText, dw.RuneglishWordText, dw.RuneWordText, dw.GemSum)
}

// DeleteAllDictionaryWords deletes all rows from the DictionaryWord table
func DeleteAllDictionaryWords(db *gorm.DB) error {
	result := db.Exec("DELETE FROM public.dictionary_words")
	return result.Error
}

// InsertDictionaryWord inserts a new DictionaryWord into the database
func InsertDictionaryWord(db *gorm.DB, word DictionaryWord) error {
	result := db.Create(&word)
	return result.Error
}

// GetWordsByLength retrieves words based on their length and text type
func GetWordsByLength(db *gorm.DB, length int, textType runer.TextType) ([]DictionaryWord, error) {
	var words []DictionaryWord
	var err error

	switch textType {
	case runer.Latin:
		err = db.Model(&DictionaryWord{}).Where("dict_word_length = ?", length).Find(&words).Error
	case runer.Runeglish:
		err = db.Model(&DictionaryWord{}).Where("dict_runeglish_length = ?", length).Find(&words).Error
	case runer.Runes:
		err = db.Model(&DictionaryWord{}).Where("dict_rune_length = ?", length).Find(&words).Error
	}

	return words, err
}

// GetWordsByGemSum retrieves words based on their gem sum
func GetWordsByGemSum(db *gorm.DB, gemSum int64) ([]DictionaryWord, error) {
	var words []DictionaryWord
	err := db.Model(&DictionaryWord{}).Where("gem_sum = ?", gemSum).Find(&words).Error
	return words, err
}

// GetWordsByPattern retrieves words based on their rune pattern
func GetWordsByPattern(db *gorm.DB, pattern string) ([]DictionaryWord, error) {
	var words []DictionaryWord
	err := db.Model(&DictionaryWord{}).Where("rune_pattern = ?", pattern).Find(&words).Error
	return words, err
}
