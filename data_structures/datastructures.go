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
	allNumbers = append(allNumbers, firstThree...) // no such method in array

	fmt.Println("The all numbers slice is as : ", allNumbers);

	fruits := []string{"apple", "banana", "mango"}
	fmt.Printf("these are my fruits: %v\n", fruits);

	// appending single
	fruits = append(fruits, "kiwi")
	fmt.Printf("these are my fruits with kiwi: %v\n", fruits);

	// appending multiple
	fruits = append(fruits, "litchi", "grapes")
	fmt.Printf("these are my fruits with litchi and grapes: %v\n", fruits);

	moreFruits := []string{"strawberry", "tomato"}
	fruits = append(fruits, moreFruits...)
	fmt.Printf("these are my fruits with more fruits: %v\n", fruits);


}
