package main

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"liberdatabase"
	"math/big"
)

// factorize returns the prime factors of a given big integer.
func factorize(db *gorm.DB, mainId string, n *big.Int, lastSeq int64) bool {
	counter := big.NewInt(2)
	zero := big.NewInt(0)
	number := new(big.Int).Set(n)

	if lastSeq > 0 {
		lastRecord := liberdatabase.GetMaxSeqNumberByMainID(db, mainId)
		liberdatabase.RemoveFactorByID(db, lastRecord.ID)
	}

	// Check if n is divisible by x
	for counter.Cmp(number) <= 0 {
		stringVal := counter.String()
		if len(stringVal) > 1 {
			if stringVal[len(stringVal)-1:] == "5" ||
				stringVal[len(stringVal)-1:] == "0" ||
				stringVal[len(stringVal)-1:] == "2" ||
				stringVal[len(stringVal)-1:] == "4" ||
				stringVal[len(stringVal)-1:] == "6" ||
				stringVal[len(stringVal)-1:] == "8" {
				// Skip even numbers, 5, and 0
				counter.Add(counter, big.NewInt(1))
				continue
			}
		}

		if new(big.Int).Mod(number, counter).Cmp(zero) == 0 {
			number = n.Div(number, counter)

			// Insert the counter factor into the database
			lastSeq++
			counterFactor := liberdatabase.Factor{
				ID:        uuid.New().String(),
				Factor:    counter.String(),
				MainId:    mainId,
				SeqNumber: lastSeq,
			}

			liberdatabase.InsertFactor(db, counterFactor)

			// Insert the number factor into the database
			lastSeq++
			numberFactor := liberdatabase.Factor{
				ID:        uuid.New().String(),
				Factor:    number.String(),
				MainId:    mainId,
				SeqNumber: lastSeq,
			}

			liberdatabase.InsertFactor(db, numberFactor)
			break
		} else {
			counter.Add(counter, big.NewInt(1))
		}
	}

	// Check if all factors are prime
	areAllPrime := true
	lastSeqNumber := int64(0)

	// Loop to get factors until nil is returned
	for {
		factor := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
		if factor == nil {
			break
		}

		// Update the last sequence number
		lastSeqNumber = factor.SeqNumber

		factorValue := new(big.Int)
		factorValue, ok := factorValue.SetString(factor.Factor, 10)
		if !ok {
			fmt.Printf("Error converting factor to *big.Int: %v\n", factor.Factor)
		}

		if !factorValue.ProbablyPrime(20) {
			areAllPrime = false
			break
		}
	}

	if areAllPrime {
		return true
	} else {
		return factorize(db, mainId, number, lastSeq)
	}
}
