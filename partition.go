// partition package contains functions to solve partition problem.
// https://www.geeksforgeeks.org/partition-problem-dp-18/
// https://www.geeksforgeeks.org/print-equal-sum-sets-array-partition-problem/
// https://www.geeksforgeeks.org/print-equal-sum-sets-array-partition-problem-set-2/

package partition

import (
	"fmt"
	"math"
	"sort"
)

func sumInt(arr []int) int {
	var sum int
	for _, val := range arr {
		sum += val
	}
	return sum
}

func sumFloat(arr []float64) float64 {
	var sum float64
	for _, val := range arr {
		sum += val
	}
	return sum
}

func absInt(value int) int {
	if value >= 0 {
		return value
	} else {
		return -value
	}
}

// isSubsetSum is utility function that returns true if there is a subset of arr[] with sum equal to given sum.
func isSubsetSum(arr []int, n, sum int) bool {
	// Base Cases
	if sum == 0 {
		return true
	}
	if n == 0 && sum != 0 {
		return false
	}

	// If last element is greater than sum, then ignore it.
	if arr[n-1] > sum {
		return isSubsetSum(arr, n-1, sum)
	}

	/* else, check if sum can be obtained by any of the following:
	(a) including the last element
	(b) excluding the last element
	*/
	return isSubsetSum(arr, n-1, sum) || isSubsetSum(arr, n-1, sum-arr[n-1])
}

// FindPartitionRecursive returns true if arr[] can be partitioned in two subsets of equal sum, otherwise false. Uses recursive function isSubsetSum.
func FindPartitionRecursive(arr []int) bool {
	// Calculate sum of the elements in array
	sum := sumInt(arr)

	// If sum is odd, there cannot be two subsets with equal sum
	if sum%2 != 0 {
		return false
	}

	// Find if there is subset with sum equal to half of total sum.
	return isSubsetSum(arr, len(arr), sum/2)
}

// FindPartitionDynamic returns true if arr[] can be partitioned in two subsets of equal sum, otherwise false. Uses dynamic programming approach.
func FindPartitionDynamic(arr []int) bool {
	// Calculate sum of all elements.
	sum := sumInt(arr)

	if sum%2 != 0 {
		return false
	}

	rows := sum/2 + 1
	columns := len(arr) + 1

	// Initialize partition table.
	part := make([][]bool, rows)
	for i := range part {
		part[i] = make([]bool, columns)
	}

	// initialize top row as true
	for i := range part[0] {
		part[0][i] = true
	}

	// We should not initialize leftmost column as false because part is already initialized by falses.

	// Fill the partition table in bottom up manner.
	for i := 1; i < rows; i++ {
		for j := 1; j < columns; j++ {
			part[i][j] = part[i][j-1]
			if i >= arr[j-1] {
				part[i][j] = part[i][j] || part[i-arr[j-1]][j-1]
			}
		}
	}

	/* // Uncomment to print table
	for i := range part {
		for j := range part[i] {
			fmt.Printf("%6v", part[i][j])
		}
		fmt.Println()
	}
	fmt.Println() */

	return part[rows-1][columns-1]
}

// Greedy makes an attempt to partition arr into two sets of equal or closest sum.
func Greedy(arr []int) ([]int, []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(arr)))
	set1 := make([]int, 0, len(arr))
	set2 := make([]int, 0, len(arr))
	var sum1, sum2 int

	for _, val := range arr {
		if sum1 < sum2 {
			set1 = append(set1, val)
			sum1 += val
		} else {
			set2 = append(set2, val)
			sum2 += val
		}
	}
	return set1, set2
}

// findSets finds the sets of the array which have equal sum.
func findSets(arr, set1, set2 []int, sum1, sum2, pos int) (bool, []int, []int) {
	// If entire array is traversed, compare both the sums.
	if pos == len(arr) {
		// If sums are equal print both sets and return true to show sets are found.
		if sum1 == sum2 {
			return true, set1, set2
		} else {
			// If sums are not equal then return sets are not found.
			return false, nil, nil
		}
	}

	// Add current element to set1.
	set1 = append(set1, arr[pos])

	// Recursive call after adding current element to set1.
	res, resSet1, resSet2 := findSets(arr, set1, set2, sum1+arr[pos], sum2, pos+1)

	// If this inclusion results in equal sum sets partition then return true to show desired sets are found.
	if res {
		return res, resSet1, resSet2
	}

	// If not then backtrack by removing current element from set1 and include it in set2.
	set1 = set1[:len(set1)-1]
	set2 = append(set2, arr[pos])

	// Recursive call after including current element to set2.
	return findSets(arr, set1, set2, sum1, sum2+arr[pos], pos+1)
}

// FindSets finds the sets of the array which have equal sum.
func FindSetsRecursive(arr []int) (bool, []int, []int) {
	// If sum of entire arr is odd then array cannot be partitioned.
	if sumInt(arr)%2 != 0 {
		return false, nil, nil
	}

	initialSet1 := make([]int, 0, len(arr))
	initialSet2 := make([]int, 0, len(arr))
	return findSets(arr, initialSet1, initialSet2, 0, 0, 0)
}

// FindSetsDynamic tries to return equal sum sets of array.
func FindSetsDynamic(arr []int) (bool, []int, []int) {
	sumArray := sumInt(arr)
	n := len(arr)

	// Check sum is even or odd. If odd then array cannot be partitioned.
	if sumArray&1 == 1 {
		return false, nil, nil
	}

	// Divide sum by 2 to find sum of two possible subsets.
	k := sumArray >> 1

	// Boolean DP table to store result of states.
	// dp[i][j] = true if there is a subset of elements in first i element of array that has sum equal to j.
	dp := make([][]bool, n+1)
	for i := range dp {
		dp[i] = make([]bool, k+1)
	}

	// If number of elements are zero, then no sum can be obtained.
	for i := 1; i <= k; i++ {
		dp[0][i] = false
	}

	// Sum 0 can be obtained by not selecting any element.
	for i := 0; i <= n; i++ {
		dp[i][0] = true
	}

	// Fill the DP table in bottom up manner.
	for i := 1; i <= n; i++ {
		for currSum := 1; currSum <= k; currSum++ {
			// Excluding current element.
			dp[i][currSum] = dp[i-1][currSum]

			// Including current element.
			if arr[i-1] <= currSum {
				dp[i][currSum] = dp[i][currSum] || dp[i-1][currSum-arr[i-1]]
			}
		}
	}

	// Required sets set1 and set2.
	set1 := make([]int, n)
	set2 := make([]int, n)

	// If partition is not possible return false.
	if !dp[n][k] {
		return false, nil, nil
	}

	// Start from last element in dp table.
	i := n
	currSum := k

	for i > 0 && currSum >= 0 {
		// If current element does not contribute to k, then it belongs to set 2.
		if dp[i-1][currSum] {
			i--
			set2 = append(set2, arr[i])
		} else if dp[i-1][currSum-arr[i-1]] {
			// If current element contribute to k then it belongs to set 1.
			i--
			currSum -= arr[i]
			set1 = append(set1, arr[i])
		}
	}

	return true, set1, set2
}

type setPair struct {
	set1 []int
	sum1 int
	set2 []int
	sum2 int
}

func (sp *setPair) sumDiff() int {
	return absInt(sp.sum1 - sp.sum2)
}

type setPairFloat struct {
	set1 []float64
	sum1 float64
	set2 []float64
	sum2 float64
}

func (spf *setPairFloat) sumDiff() float64 {
	return math.Abs(spf.sum1 - spf.sum2)
}

func getMinSetPair(sp1, sp2 setPair) setPair {
	if sp1.sumDiff() <= sp2.sumDiff() {
		return sp1
	} else {
		return sp2
	}
}

func getMinSetPairFloat(sp1, sp2 setPairFloat) setPairFloat {
	if sp1.sumDiff() <= sp2.sumDiff() {
		return sp1
	} else {
		return sp2
	}
}

// findMinSetPair partition array into two sets such that the difference of set sums is minimum. Uses recursive approach.
func findMinSetPair(arr []int, sp setPair, pos int) setPair {
	// If entire array is traversed, return result.
	if pos == len(arr) {
		return sp
	}

	// For every item arr[pos] we have two choices:
	// 1. We include it into the first set;
	// 2. We include it into the second set.
	// Return minimum of two results.
	return getMinSetPair(
		findMinSetPair(arr, setPair{append(sp.set1, arr[pos]), sp.sum1 + arr[pos], append([]int{}, sp.set2...), sp.sum2}, pos+1),
		findMinSetPair(arr, setPair{append([]int{}, sp.set1...), sp.sum1, append(sp.set2, arr[pos]), sp.sum2 + arr[pos]}, pos+1),
	)
}

// FindSetsWithMinSumDifference is the wrapper over findMinSetPair.
func FindSetsWithMinSumDifference(arr []int) ([]int, []int, int) {
	minSetPair := findMinSetPair(arr, setPair{[]int{}, 0, []int{}, 0}, 0)
	return minSetPair.set1, minSetPair.set2, minSetPair.sumDiff()
}

// findMinEqualSetPair partition array into two equal in number sets such the difference of set sumsis minimum. Uses recursive approach.
func findMinEqualSetPair(arr []int, sp setPair, pos int) setPair {
	// If size of whether set in setPair is equal the half of arr size, put the rest elements to another set and return result.
	if len(sp.set1) == len(arr)/2 || len(sp.set2) == len(arr)/2 {
		if len(sp.set1) == len(arr)/2 {
			sp.set2 = append(sp.set2, arr[pos:]...)
			sp.sum2 += sumInt(arr[pos:])
		} else {
			sp.set1 = append(sp.set1, arr[pos:]...)
			sp.sum1 += sumInt(arr[pos:])
		}
		return sp
	}

	// For every item arr[pos] we have two choices:
	// 1. We include it into the first set;
	// 2. We include it into the second set.
	// Return minimum of two results.
	return getMinSetPair(
		findMinEqualSetPair(arr, setPair{append(sp.set1, arr[pos]), sp.sum1 + arr[pos], append([]int{}, sp.set2...), sp.sum2}, pos+1),
		findMinEqualSetPair(arr, setPair{append([]int{}, sp.set1...), sp.sum1, append(sp.set2, arr[pos]), sp.sum2 + arr[pos]}, pos+1),
	)
}

// FindEqualSetsWithMinSumDifference is the wrapper over findMinEqualSetPair.
func FindEqualSetsWithMinSumDifference(arr []int) ([]int, []int, int, error) {
	// If size of arr is odd, we can't partition it on two equal sets.
	if len(arr)&1 != 0 {
		return nil, nil, 0, fmt.Errorf("array size is odd. Can't find equal number sets.\n")
	}
	minEqualSetPair := findMinEqualSetPair(arr, setPair{[]int{}, 0, []int{}, 0}, 0)
	return minEqualSetPair.set1, minEqualSetPair.set2, minEqualSetPair.sumDiff(), nil
}

// findMinEqualSetPairFloat partition array into two equal in number sets such the difference of set sums is minimum. Uses recursive approach.
func findMinEqualSetPairFloat(arr []float64, sp setPairFloat, pos int) setPairFloat {
	// If size of whether set in setPair is equal the half of arr size, put the rest elements to another set and return result.
	if len(sp.set1) == len(arr)/2 || len(sp.set2) == len(arr)/2 {
		if len(sp.set1) == len(arr)/2 {
			sp.set2 = append(sp.set2, arr[pos:]...)
			sp.sum2 += sumFloat(arr[pos:])
		} else {
			sp.set1 = append(sp.set1, arr[pos:]...)
			sp.sum1 += sumFloat(arr[pos:])
		}
		return sp
	}

	// For every item arr[pos] we have two choices:
	// 1. We include it into the first set;
	// 2. We include it into the second set.
	// Return minimum of two results.
	return getMinSetPairFloat(
		findMinEqualSetPairFloat(arr, setPairFloat{append(sp.set1, arr[pos]), sp.sum1 + arr[pos], append([]float64{}, sp.set2...), sp.sum2}, pos+1),
		findMinEqualSetPairFloat(arr, setPairFloat{append([]float64{}, sp.set1...), sp.sum1, append(sp.set2, arr[pos]), sp.sum2 + arr[pos]}, pos+1),
	)
}

// FindEqualSetsWithMinSumDifferenceFloat is the wrapper over findMinEqualSetPairFloat.
func FindEqualSetsWithMinSumDifferenceFloat(arr []float64) ([]float64, []float64, float64, error) {
	// If size of arr is odd, we can't partition it on two equal sets.
	if len(arr)&1 != 0 {
		return nil, nil, 0, fmt.Errorf("array size is odd. Can't find equal number sets.\n")
	}
	minEqualSetPair := findMinEqualSetPairFloat(arr, setPairFloat{[]float64{}, 0, []float64{}, 0}, 0)
	return minEqualSetPair.set1, minEqualSetPair.set2, minEqualSetPair.sumDiff(), nil
}
