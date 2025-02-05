package sequences

import (
	"config"
	"fmt"
	"math/big"
	"sync"
)

// IsPrime checks if a number is prime.
func IsPrime(number *big.Int) bool {
	if number.Cmp(big.NewInt(2)) < 0 {
		return false
	}
	if number.Cmp(big.NewInt(2)) == 0 {
		return true
	}
	if new(big.Int).Mod(number, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		return false
	}

	sqrt := new(big.Int).Sqrt(number)
	for i := big.NewInt(3); i.Cmp(sqrt) <= 0; i.Add(i, big.NewInt(2)) {
		if new(big.Int).Mod(number, i).Cmp(big.NewInt(0)) == 0 {
			return false
		}
	}

	return true
}

// YieldPrimesAsc yields prime numbers in descending order up to the given number.
func YieldPrimesAsc(maxNumber *big.Int) <-chan *big.Int {
	one := big.NewInt(1)
	ch := make(chan *big.Int)
	var wg sync.WaitGroup

	// Load worker count from config
	cfg, err := config.LoadConfig()
	workerCount := 4 // Default worker count
	if err != nil {
		fmt.Printf("Error loading config: %v\nUsing default worker count: %d\n", err)
	} else {
		workerCount = cfg.NumWorkers / 2
	}

	// Calculate the range size for each worker
	rangeSize := new(big.Int).Div(maxNumber, big.NewInt(int64(workerCount)))

	// Start worker goroutines
	for i := 0; i < workerCount; i++ {
		start := new(big.Int).Mul(rangeSize, big.NewInt(int64(i)))
		end := new(big.Int).Add(start, rangeSize)
		if i == workerCount-1 {
			end = maxNumber
		}

		wg.Add(1)
		go func(start, end *big.Int, isDecrement bool) {
			defer wg.Done()
			// if increment is true, the worker will yield prime numbers in ascending order
			counter := new(big.Int).Set(start)

			if isDecrement {
				counter = new(big.Int).Set(end)
			}

			for {
				if isDecrement {
					if counter.Cmp(start) <= 0 {
						break
					}
					if counter.ProbablyPrime(20) {
						ch <- new(big.Int).Set(counter)
					}
					counter.Sub(counter, one)
				} else {
					if counter.Cmp(end) >= 0 {
						break
					}
					if counter.ProbablyPrime(20) {
						ch <- new(big.Int).Set(counter)
					}
					counter.Add(counter, one)
				}
			}
		}(start, end, i < workerCount/2)
	}

	// Close the channel once all workers are done
	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch
}

// GetPrimeSequence generates the prime sequence.
func GetPrimeSequence(maxNumber *big.Int, isPositional bool) (*NumericSequence, error) {
	numericSequence := &NumericSequence{Name: "Prime", Number: new(big.Int).Set(maxNumber)}
	numberToCalculate := new(big.Int).Set(maxNumber)
	if isPositional {
		numberToCalculate = new(big.Int).SetUint64(^uint64(0)) // Max uint64 value
	}
	counter := big.NewInt(0)

	for i := big.NewInt(0); i.Cmp(numberToCalculate) <= 0; i.Add(i, big.NewInt(1)) {
		if IsPrime(i) {
			if !isPositional {
				numericSequence.Sequence = append(numericSequence.Sequence, new(big.Int).Set(i))
			} else {
				if counter.Cmp(maxNumber) == 0 {
					numericSequence.Sequence = append(numericSequence.Sequence, new(big.Int).Set(i))
					break
				}
			}
			counter.Add(counter, big.NewInt(1))
		}
	}

	return numericSequence, nil
}

// GetFibonacciPrimeSequence generates the Fibonacci prime sequence.
func GetFibonacciPrimeSequence(maxNumber *big.Int, isPositional bool) (*NumericSequence, error) {
	numericSequence := &NumericSequence{Name: "Fibonacci Prime", Number: new(big.Int).Set(maxNumber)}
	numberToCalculate := new(big.Int).Set(maxNumber)
	if isPositional {
		numberToCalculate = new(big.Int).SetUint64(^uint64(0)) // Max uint64 value
	}

	a, b, c := big.NewInt(0), big.NewInt(1), big.NewInt(0)
	counter := big.NewInt(0)

	for c.Cmp(numberToCalculate) <= 0 {
		c.Add(a, b)
		a.Set(b)
		b.Set(c)

		if c.Cmp(numberToCalculate) <= 0 && IsPrime(c) {
			if !isPositional {
				numericSequence.Sequence = append(numericSequence.Sequence, new(big.Int).Set(c))
			} else {
				if counter.Cmp(maxNumber) == 0 {
					numericSequence.Sequence = append(numericSequence.Sequence, new(big.Int).Set(c))
					break
				}
			}
			counter.Add(counter, big.NewInt(1))
		}
	}

	return numericSequence, nil
}
