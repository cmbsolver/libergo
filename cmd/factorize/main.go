package main

import (
	"flag"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"liberdatabase"
	"math/big"
	"os"
	"runtime"
	"sequences"
	"sort"
	"strings"
	"sync"
	"time"
	"titler"
)

type NumberToCheck struct {
	Number  string
	Counter string
}

var status big.Int

func main() {
	titler.PrintTitle("Factorize")
	startTime := time.Now() // Record the start time

	// Timer to write out the status every minute
	ticker := time.NewTicker(1 * time.Minute)
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-ticker.C:
				fmt.Printf("Status update: %s (bits %d)\n", status.String(), status.BitLen())
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()
	// Make sure to stop the ticker when the program ends
	defer func() {
		done <- true
	}()

	// Parse command-line flags
	flag.Parse()

	// Check if the number is provided as an argument
	if len(flag.Args()) < 1 {
		fmt.Println("Please provide a number to be factorized as an argument.")
		os.Exit(1)
	}

	// Read the input number
	numberStr := flag.Arg(0)

	// Convert input to bigint
	number := new(big.Int)
	number, ok := number.SetString(numberStr, 10)
	if !ok {
		fmt.Println("Invalid number format.")
		os.Exit(1)
	}

	if number.Cmp(big.NewInt(1)) == -1 || number.Cmp(big.NewInt(1)) == 0 {
		fmt.Printf("%s : 1,%s\n", numberStr, numberStr)
		return
	}

	if sequences.IsPrime(number) {
		// You don't need to factorize a prime number
		fmt.Printf("%s : 1,%s\n", numberStr, numberStr)
		return
	}

	// Establish database connection
	cError := liberdatabase.InitSQLiteTables()
	if cError != nil {
		fmt.Printf("Error creating database tables: %v\n", cError)
		os.Exit(1)
	}

	db, connError := liberdatabase.InitSQLiteConnection()
	if connError != nil {
		fmt.Printf("Error initializing database connection: %v\n", connError)
		os.Exit(1)
	}

	// The mainId is the number being factorized
	mainId := uuid.New().String()

	fmt.Printf("Factorizing %s (%d bits)\n", number.String(), number.BitLen())

	// Perform factorization
	factorize(db, mainId, number, 0)

	// Output prime factors
	output := strings.Builder{}
	firstTime := true

	// Initialize the last sequence number
	var lastSeqNumber int64 = 0

	// Loop to get factors until nil is returned
	for {
		factor := liberdatabase.GetFactorsByMainID(db, mainId, lastSeqNumber)
		if factor == nil {
			break
		}

		// Update the last sequence number
		lastSeqNumber = factor.SeqNumber

		if !firstTime {
			output.WriteString(",")
		}

		// Append factor to output
		output.WriteString(factor.Factor)

		firstTime = false
	}

	fmt.Println(numberStr, ":", output.String())
	writeOutputToFile(fmt.Sprintf("%s : %s", numberStr, output.String()))

	liberdatabase.RemoveFactorsByMainID(db, mainId)

	endTime := time.Now()                        // Record the end time
	duration := endTime.Sub(startTime)           // Calculate the duration
	fmt.Printf("Execution time: %v\n", duration) // Print the log message
}

// factorize returns the prime factors of a given big integer.
func factorize(db *gorm.DB, mainId string, n *big.Int, lastSeq int64) bool {
	zero := big.NewInt(0)
	number := new(big.Int).Set(n)
	var modNumberArray []big.Int
	processedCounter := big.NewInt(0)

	fmt.Printf("- Step - Factoring %s (%d bits)\n", number.String(), number.BitLen())

	if lastSeq > 0 {
		lastRecord := liberdatabase.GetMaxSeqNumberByMainID(db, mainId)
		liberdatabase.RemoveFactorByID(db, lastRecord.ID)
	}

	if number.ProbablyPrime(20) {
		fmt.Printf("-%s is prime\n", number.String())
		modNumberArray = append(modNumberArray, *number)
	} else {

		// We're going to use threads to check this out
		checkChannel := make(chan NumberToCheck)
		var wg sync.WaitGroup
		numProcessors := runtime.NumCPU()
		waits := 0
		if number.Cmp(big.NewInt(int64(numProcessors*4))) > 0 {
			waits = numProcessors * 4
		} else {
			waits = numProcessors
		}

		// Start worker goroutines
		for i := 0; i < numProcessors*2; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for checkValue := range checkChannel {
					myBigNumber, _ := new(big.Int).SetString(checkValue.Number, 10)
					myBigCounter, _ := new(big.Int).SetString(checkValue.Counter, 10)
					if new(big.Int).Mod(myBigNumber, myBigCounter).Cmp(zero) == 0 {
						modNumberArray = append(modNumberArray, *myBigCounter)
						break
					}
				}
			}()
		}

		go func() {
			one := big.NewInt(1)
			two := big.NewInt(2)
			squareRoot := new(big.Int).Set(n)
			squareRoot = squareRoot.Sqrt(squareRoot)

			fmt.Printf("Square Root: %s (bits %d)\n", squareRoot.String(), squareRoot.BitLen())

			myCounter := new(big.Int).Set(two)

			for myCounter.Cmp(squareRoot) <= 0 {
				status.Set(myCounter)

				checkChannel <- NumberToCheck{
					Number:  number.String(),
					Counter: myCounter.String(),
				}

				myCounter.Add(myCounter, one)

				if len(modNumberArray) > 0 && processedCounter.Cmp(big.NewInt(int64(waits))) > 0 {
					break
				}

				processedCounter.Add(processedCounter, one)
			}

			myCounter = new(big.Int).Set(number)
			myCounter.Sub(myCounter, one)

			for myCounter.Cmp(squareRoot) >= 0 {
				status.Set(myCounter)

				checkChannel <- NumberToCheck{
					Number:  number.String(),
					Counter: myCounter.String(),
				}

				myCounter.Sub(myCounter, one)

				if len(modNumberArray) > 0 && processedCounter.Cmp(big.NewInt(int64(waits))) > 0 {
					break
				}

				processedCounter.Add(processedCounter, one)
			}
			close(checkChannel)
		}()

		wg.Wait()

		// grabbing the smallest factor
		modNumberArray = sortBigInts(modNumberArray)
	}

	bcounter := modNumberArray[0]

	fmt.Printf("- Factor %s found\n", bcounter.String())

	number = n.Div(number, &bcounter)

	// Insert the count factor into the database
	lastSeq++
	counterFactor := liberdatabase.Factor{
		ID:        uuid.New().String(),
		Factor:    bcounter.String(),
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

func sortBigInts(bigInts []big.Int) []big.Int {
	sort.Slice(bigInts, func(i, j int) bool {
		return bigInts[i].Cmp(&bigInts[j]) < 0
	})
	return bigInts
}

func writeOutputToFile(output string) {
	file, err := os.OpenFile("factorize_output.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer func(file *os.File) {
		closeError := file.Close()
		if closeError != nil {
			fmt.Printf("Error closing file: %v\n", closeError)
		}
	}(file)

	if _, writeError := file.WriteString(output + "\n"); writeError != nil {
		fmt.Printf("Error writing to file: %v\n", writeError)
	}
}
