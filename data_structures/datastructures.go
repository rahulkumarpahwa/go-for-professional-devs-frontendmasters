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

	// maps:

	capitalCities := map[string]string{
		"India" : "New Delhi",
		"USA" : "Washington D.C.",
		"UK" : "London",
	}

	fmt.Printf("The capital of the USA is: %v\n", capitalCities["USA"])
	fmt.Printf("The capital of the UK is: %v\n", capitalCities["UK"])
	fmt.Printf("The capital of Germany is: %v\n", capitalCities["Germany"]) // if key is not present it will return zero value of the type which is empty string in this case

	// checking if a key exists in the map
	capital, exists := capitalCities["Germany"]
	if exists {
		fmt.Printf("The capital of Germany is: %v\n", capital)
	} else {
		fmt.Println("Germany is not in the map.")
	}

	// inbuilt function to get the length of the map
	fmt.Printf("The number of countries in the map is: %v\n", len(capitalCities))
	// deleting a key from the map
	delete(capitalCities, "UK") // deleting a key from the map
	fmt.Printf("After deleting UK, the map is: %v\n", capitalCities);
}