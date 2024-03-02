package grid

import (
	"fmt"
	"math/rand"
	"sort"
)

type SudokuMap map[int]int

func createMap(inputSize, outputSize, seed int) (SudokuMap, error) {
	if inputSize <= 0 || outputSize <= 0 {
		return nil, fmt.Errorf("input and output sizes must be positive, got: %d, %d", inputSize, outputSize)
	}
	sm := make(SudokuMap)
	for i := range inputSize {
		// Need to cache the map with a hash at some point for efficiency...
		sm[i] = rand.New(rand.NewSource(int64(seed) * int64(i))).Intn(outputSize)
	}
	return sm, nil
}

func (sm SudokuMap) Preimage() map[int][]int {
	preimage := make(map[int][]int)
	for k, v := range sm {
		preimage[v] = append(preimage[v], k)
	}
	for k := range preimage {
		sort.Ints(preimage[k])
	}
	return preimage
}

func Identity(fullDimension int) (SudokuMap, error) {
	if fullDimension <= 0 {
		return nil, fmt.Errorf("size must be positive, got: %d", fullDimension)
	}
	identity := make(SudokuMap)
	for i := range fullDimension {
		identity[i] = i
	}
	return identity, nil
}

// func (sm SudokuMap) apply(input []int) []int {
// 	output := make([]int, len(input))
// 	for i, v := range input {
// 		output[i] = sm[v]
// 	}
// 	return output
// }

// func createInjective(inputSize, outputSize int) (SudokuMap, error) {
// 	return make(SudokuMap), nil
// }

// func createNonInjective(inputSize, outputSize int) (SudokuMap, error) {
// 	return make(SudokuMap), nil
// }

// func createSurjective(inputSize, outputSize int) (SudokuMap, error) {
// 	return make(SudokuMap), nil
// }

// func createNonSurjective(inputSize, outputSize int) (SudokuMap, error) {
// 	return make(SudokuMap), nil
// }

// // isomorphisms are maps that preserve neighbors.
// func createIsomorphism(inputSize int) (SudokuMap, error) {
// 	return make(SudokuMap), nil
// }
