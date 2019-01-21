package partition

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"testing"
)

const (
	arraysNumber = 1000
	arraySize    = 12
	maxInt       = 100
	seed         = 42
)

var (
	arrays                [arraysNumber][arraySize]int
	partitionableArrays   [][arraySize]int
	unpartitionableArrays [][arraySize]int
)

func TestMain(m *testing.M) {
	rand.Seed(seed) // The same sets every time

	for i, array := range arrays {
		for j := range array {
			array[j] = rand.Intn(maxInt)
		}
		arrays[i] = array
	}

	os.Exit(m.Run())
}

func findPartition(findPartitionFunc func([]int) bool) {
	// Clear arrays set
	partitionableArrays = partitionableArrays[:0]
	unpartitionableArrays = unpartitionableArrays[:0]

	for _, array := range arrays {
		if findPartitionFunc(array[:]) {
			partitionableArrays = append(partitionableArrays, array)
		} else {
			unpartitionableArrays = append(unpartitionableArrays, array)
		}
	}
}

func testFindPartition(t *testing.T, findPartitionFunc func([]int) bool) {
	findPartition(findPartitionFunc)

	if len(partitionableArrays) != 502 || len(unpartitionableArrays) != 498 {
		t.Errorf("Wrong results: expected 502 partitionable arrays and 498 unpartitionable, but got %d and %d respectively.", len(partitionableArrays), len(unpartitionableArrays))
	}
}

func TestFindPartitionRecursive(t *testing.T) {
	testFindPartition(t, FindPartitionRecursive)
}

func TestFindPartitionDynamic(t *testing.T) {
	testFindPartition(t, FindPartitionDynamic)
}

func TestGreedy(t *testing.T) {
	if len(partitionableArrays) == 0 && len(unpartitionableArrays) == 0 {
		findPartition(FindPartitionRecursive)
	}

	results := map[bool]int{true: 0, false: 0}

	for _, array := range partitionableArrays {
		set1, set2 := Greedy(array[:])
		sum1 := sumInt(set1)
		sum2 := sumInt(set2)
		results[sum1 == sum2]++
		// Greedy approach must give 7/6-approximation: https://en.wikipedia.org/wiki/Partition_problem#The_greedy_algorithm
		if sum1 != sum2 {
			var maxSum int
			if sum1 > sum2 {
				maxSum = sum1
			} else {
				maxSum = sum2
			}
			averageSum := (sum1 + sum2) / 2
			if float64(maxSum)/float64(averageSum) > 7.0/6.0 {
				t.Errorf("Wrong partition of array %v on %v and %v by greedy algorythm. Sums %v %v have too big difference (maxSum = %v, average sum = %v).", array, set1, set2, sum1, sum2, maxSum, averageSum)
			}
		}
	}

	if results[true] != 92 || results[false] != 410 {
		t.Errorf("Wrong results: expected 92 true and 410 false, but got %d and %d respectively.", results[true], results[false])
	}
}

func testFindSets(t *testing.T, findSetsFunc func([]int) (bool, []int, []int)) {
	results := map[bool]int{true: 0, false: 0}

	for _, array := range arrays {
		res, set1, set2 := findSetsFunc(array[:])
		results[res]++

		if res {
			if sumInt(set1) != sumInt(set2) {
				t.Errorf("Wrong partition of array %v on %v and %v: sums are not equal (%v != %v)", array, set1, set2, sumInt(set1), sumInt(set2))
			}
		}
	}

	if results[true] != 502 || results[false] != 498 {
		t.Errorf("Wrong results: expected 502 partitionable arrays and 498 unpartitionable, but got %d and %d respectively.", results[true], results[false])
	}
}

func TestFindSetsRecursive(t *testing.T) {
	testFindSets(t, FindSetsRecursive)
}

func TestFindSetsDynamic(t *testing.T) {
	testFindSets(t, FindSetsDynamic)
}

func compareArrayAndSets(t *testing.T, array [arraySize]int, set1, set2 []int) {
	var sliceFromArray = make([]int, 0, arraySize)
	var sliceFromSets = make([]int, 0, arraySize)

	sliceFromArray = append(sliceFromArray, array[:]...)
	sliceFromSets = append(sliceFromSets, set1...)
	sliceFromSets = append(sliceFromSets, set2...)

	sort.Ints(sliceFromArray)
	sort.Ints(sliceFromSets)

	for len(sliceFromArray) != len(sliceFromSets) {
		t.Errorf("Mismatch of lengths of array %v and sets %v and %v: %d != %d + %d\n", array, set1, set2, len(array), len(set1), len(set2))
	}
	for i := range sliceFromArray {
		if sliceFromArray[i] != sliceFromSets[i] {
			t.Errorf("Mismatch in contents of array %v and sets %v and %v: e.g. %d\n", array, set1, set2, sliceFromArray[i])
			break
		}
	}
}

func compareArrayAndSetsFloat(t *testing.T, array [arraySize]float64, set1, set2 []float64) {
	var sliceFromArray = make([]float64, 0, arraySize)
	var sliceFromSets = make([]float64, 0, arraySize)

	sliceFromArray = append(sliceFromArray, array[:]...)
	sliceFromSets = append(sliceFromSets, set1...)
	sliceFromSets = append(sliceFromSets, set2...)

	sort.Float64s(sliceFromArray)
	sort.Float64s(sliceFromSets)

	for len(sliceFromArray) != len(sliceFromSets) {
		t.Errorf("Mismatch of lengths of array %v and sets %v and %v: %d != %d + %d\n", array, set1, set2, len(array), len(set1), len(set2))
	}
	for i := range sliceFromArray {
		if sliceFromArray[i] != sliceFromSets[i] {
			t.Errorf("Mismatch in contents of array %v and sets %v and %v: e.g. %f\n", array, set1, set2, sliceFromArray[i])
			break
		}
	}
}

func TestFindSetsWithMinSumDifference(t *testing.T) {
	results := make(map[int]int)

	for _, array := range arrays {
		set1, set2, sumDiff := FindSetsWithMinSumDifference(array[:])
		compareArrayAndSets(t, array, set1, set2)
		results[sumDiff]++

		if absInt(sumInt(set1)-sumInt(set2)) != sumDiff {
			t.Errorf("Wrong partition of array %v on %v and %v: abs(%v - %v) != %v\n", array, set1, set2, sumInt(set1), sumInt(set2), sumDiff)
		}
	}

	validResults := map[int]int{0: 502, 1: 493, 2: 3, 3: 1, 4: 1}
	for k, v := range results {
		if validResults[k] != v {
			t.Errorf("Wrong results: expected map[0:502 1:493 2:3 3:1 4:1], got %v\n", results)
			break
		}
	}
	fmt.Println(results)
}

func TestFindEqualSetsWithMinSumDifference(t *testing.T) {
	results := make(map[int]int)

	for _, array := range arrays {
		set1, set2, sumDiff, err := FindEqualSetsWithMinSumDifference(array[:])
		if err != nil {
			t.Errorf("array %v: %v", array, err)
		}
		if len(set1) != len(set2) {
			t.Errorf("Sizes of sets %v and %v are not equal: %d != %d\n", set1, set2, len(set1), len(set2))
		}
		compareArrayAndSets(t, array, set1, set2)
		results[sumDiff]++

		if absInt(sumInt(set1)-sumInt(set2)) != sumDiff {
			t.Errorf("Wrong partition of array %v on %v and %v: abs(%v - %v) != %v\n", array, set1, set2, sumInt(set1), sumInt(set2), sumDiff)
		}
	}

	validResults := map[int]int{0: 476, 1: 490, 2: 27, 3: 4, 4: 1, 6: 1, 18: 1}
	for k, v := range results {
		if validResults[k] != v {
			t.Errorf("Wrong results: expected map[0:502 1:493 2:3 3:1 4:1], got %v\n", results)
			break
		}
	}
	fmt.Println(results)
}

func TestFindEqualSetsWithMinSumDifferenceFloat(t *testing.T) {
	// Build arrays with float numbers.
	var arrays [arraysNumber][arraySize]float64
	rand.Seed(seed)
	for i, array := range arrays {
		for j := range array {
			array[j] = rand.Float64()
		}
		arrays[i] = array
	}

	results := make([]float64, 0, arraysNumber)

	for _, array := range arrays {
		set1, set2, sumDiff, err := FindEqualSetsWithMinSumDifferenceFloat(array[:])
		if err != nil {
			t.Errorf("array %v: %v", array, err)
		}
		if len(set1) != len(set2) {
			t.Errorf("Sizes of sets %v and %v are not equal: %d != %d\n", set1, set2, len(set1), len(set2))
		}
		compareArrayAndSetsFloat(t, array, set1, set2)
		results = append(results, sumDiff)

		floatError := math.Abs(sumFloat(set1)-sumFloat(set2)) - sumDiff
		diff := math.Round(floatError*1e14) / 1e14
		if diff != 0 {
			t.Errorf("Wrong partition of array %v on %v and %v: abs(%v - %v) - %v != 0 (%v)\n", array, set1, set2, sumFloat(set1), sumFloat(set2), sumDiff, diff)
		}
	}
	sort.Float64s(results)
	fmt.Printf("sumDiffs from %v to %v\n", results[0], results[len(results)-1])
}
