// partition package contains functions to solve partition problem.
// https://www.geeksforgeeks.org/partition-problem-dp-18/
// https://www.geeksforgeeks.org/print-equal-sum-sets-array-partition-problem/
// https://www.geeksforgeeks.org/print-equal-sum-sets-array-partition-problem-set-2/

package partition

func sumInt(arr []int) int {
	var sum int
	for _, val := range arr {
		sum += val
	}
	return sum
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
	sum := sum(arr)

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
