package partition

import (
	"math/rand"
	"os"
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
