package partition

import (
	"math/rand"
	"os"
	"testing"
)

var (
	arrays                [100][12]int
	partitionableArrays   [][12]int
	unpartitionableArrays [][12]int
)

func TestMain(m *testing.M) {
	rand.Seed(42) // The same sets every time

	for i, array := range arrays {
		for j := range array {
			array[j] = rand.Intn(100)
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

	if len(partitionableArrays) != 56 || len(unpartitionableArrays) != 44 {
		t.Errorf("Wrong results: expected 56 partitionable arrays and 44 unpartitionable, but got %d and %d respectively.", len(partitionableArrays), len(unpartitionableArrays))
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

	if results[true] != 15 || results[false] != 41 {
		t.Errorf("Wrong results: expected 15 true and 41 false, but got %d and %d respectively.", results[true], results[false])
	}
}

func TestFindSets(t *testing.T) {
	results := map[bool]int{true: 0, false: 0}

	for _, array := range arrays {
		res, set1, set2 := FindSets(array[:])
		results[res]++

		if res {
			if sumInt(set1) != sumInt(set2) {
				t.Errorf("Wrong partition of array %v on %v and %v: sums are not equal (%v != %v)", array, set1, set2, sumInt(set1), sumInt(set2))
			}
		}
	}

	if results[true] != 56 || results[false] != 44 {
		t.Errorf("Wrong results: expected 56 partitionable arrays and 44 unpartitionable, but got %d and %d respectively.", results[true], results[false])
	}
}
