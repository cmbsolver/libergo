package sequences

import (
	"math/big"
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
