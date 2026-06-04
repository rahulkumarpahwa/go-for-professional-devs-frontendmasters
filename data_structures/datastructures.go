package main

import "fmt"

func main() {
	// Arrays and Slices

	numbers := [5]int{10, 20, 30, 40, 50}
	fmt.Printf("this is out an array  %v\n", numbers)

	matrix := [2][3]int{
		{}, {2, 3, 4},
	}
	fmt.Printf("this is out a matrix  %v\n", matrix)

	// Slices
	allNumbers := numbers[:]
	firstThree := numbers[0:3]
	allNumbers = append(allNumbers, firstThree...)

	
}
