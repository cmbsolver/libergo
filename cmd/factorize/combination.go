package main

import (
	"config"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"liberdatabase"
	"math/big"
	"os"
	"sequences"
	"sync"
)

// findCombos finds prime combos for a given number.
func findCombos(db *pgx.Conn, mainId string, n *big.Int) bool {
	number := new(big.Int).Set(n)
	seqNumber := int64(0)
	loopCounter := int64(0)

	// Get p values
	fmt.Println("Getting possible p values")
	getPValues(mainId, number)

	// Initialize the last sequence number
	var lastSeqNumber int64 = 0

	// Loop to get factors until nil is returned
	for {
		loopCounter++
		factor, err := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
		if err != nil {
			fmt.Printf("Error getting factors: %v\n", err)
			os.Exit(1)
		}
		if factor == nil {
			break
		}

		// Update the last sequence number
		lastSeqNumber = factor.SeqNumber

		// Convert the factor to a big.Int
		prime := new(big.Int)
		if _, ok := prime.SetString(factor.Factor, 10); !ok {
			fmt.Printf("Error converting factor to big.Int: %s\n", factor.Factor)
			continue
		}

		if loopCounter == 1000000 {
			fmt.Printf("Current prime at loop %d: %s\n", loopCounter, factor.Factor)
			loopCounter = 0 // Reset loopCounter
		}

		fmt.Println("Processing factor: ", factor.Factor)

		q := new(big.Int).Div(number, prime)

		if q.ProbablyPrime(20) {
			seqNumber++

			// Insert the prime combo into the database
			combo := liberdatabase.PrimeCombo{
				ID:        uuid.New().String(),
				ValueP:    prime.String(),
				ValueQ:    q.String(),
				MainId:    mainId,
				SeqNumber: seqNumber,
			}

			fmt.Println("Found combo: ", combo.ValueP, combo.ValueQ)

			err := liberdatabase.InsertPrimeCombo(db, combo)
			if err != nil {
				fmt.Printf("Error inserting factor: %v\n", err)
				return false
			}
		}

	}

	removeErr := liberdatabase.RemoveFactorsByMainID(db, mainId)
	if removeErr != nil {
		fmt.Printf("Error removing factors: %v\n", removeErr)
	}

	return true
}

// getPValues finds p values using multiple workers.
func getPValues(mainId string, n *big.Int) {
	// Load worker count from config
	cfg, err := config.LoadConfig()
	workerCount := 4 // Default worker count
	if err != nil {
		fmt.Printf("Error loading config: %v\nUsing default worker count: %d\n", err, workerCount)
	} else {
		workerCount = cfg.NumWorkers / 2
	}

	// Create channels for distributing work and collecting results
	primeChan := make(chan *big.Int)
	resultChan := make(chan *big.Int)
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			// Each worker initializes its own database connection
			db, err := liberdatabase.InitConnection()
			if err != nil {
				fmt.Printf("Error initializing database connection: %v\n", err)
				return
			}
			defer func(db *pgx.Conn) {
				dbCloseError := liberdatabase.CloseConnection(db)
				if dbCloseError != nil {
					fmt.Printf("Error closing database connection: %v\n", dbCloseError)
				}
			}(db)

			for prime := range primeChan {
				if new(big.Int).Mod(n, prime).Cmp(big.NewInt(0)) == 0 {
					resultChan <- prime
				}
			}
		}(i)
	}

	// Start a goroutine to close the result channel once all workers are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Start a goroutine to send primes to the workers
	go func() {
		for prime := range sequences.YieldPrimesAsc(n) {
			primeChan <- prime
		}
		close(primeChan)
	}()

	seqValue := int64(0)

	// Collect results
	for prime := range resultChan {
		seqValue++
		fmt.Printf("Found prime factor: %s\n", prime.String())
		// Insert the prime into the database or perform other actions
		factor := liberdatabase.Factor{
			ID:        uuid.New().String(),
			Factor:    prime.String(),
			MainId:    mainId,
			SeqNumber: seqValue,
		}

		// Create a new database connection for inserting the factor
		db, err := liberdatabase.InitConnection()
		if err != nil {
			fmt.Printf("Error initializing database connection: %v\n", err)
			continue
		}
		err = liberdatabase.InsertFactor(db, factor)
		if err != nil {
			fmt.Printf("Error inserting factor: %v\n", err)
		}
		dbCloseError := liberdatabase.CloseConnection(db)
		if dbCloseError != nil {
			fmt.Printf("Error closing database connection: %v\n", dbCloseError)
		}
	}
}
