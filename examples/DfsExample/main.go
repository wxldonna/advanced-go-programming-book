package main

import "fmt"

func main() {
	// there are 6 boxes , each box has one item , we have 6 items
	// how many sequence can be founded

	n := 4
	boxex := make(map[int]string)
	nums := []int{}
	seen := make(map[int]bool)

	for i := 0; i < n+1; i++ {
		boxex[i] = ""
		nums = append(nums, i)

	}

	// 0,1,2
	var dfs func(step int)
	dfs = func(step int) {
		// exit condition
		if step == n+1 {
			fmt.Printf("boxes secquence is %v \n", boxex)
			return
		}
		for i := 1; i <= n; i++ {
			if !seen[i] {
				boxex[step] = fmt.Sprintf("%d", nums[i])
				seen[i] = true
				dfs(step + 1)
				seen[i] = false

			}
		}
		return
	}

	dfs(1)

	//dfs(1)
}

/*
var n = 3
var a = [4]int{}
var book = [4]int{}

func dfs(step int) {
	if step == n+1 {
		fmt.Printf("===========%v\n", a)

	}

	for i := 1; i <= n; i++ {
		if book[i] == 0 {
			a[step] = i
			book[i] = 1
			dfs(step + 1)
			book[i] = 0

		}
	}
	return
}
*/
