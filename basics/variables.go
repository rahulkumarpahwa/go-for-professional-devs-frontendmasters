package main

import "fmt"

func main() {
	var name string = "rahul"

	fmt.Printf("My name is: %s\n", name) // note: %s is for the string

	age := 30                          // walrus operator
	fmt.Printf("My age is: %d\n", age) // note: %d is for the integer

	var city string
	// here in between when the value of the city is not defined then it is empty string.
	fmt.Printf("this is the value of the city before definition : '%s'\n", city)
	city = "Seattle"
	fmt.Printf("My City is: %s\n", city)

	// multiple declaration:
	var country, continent string = "USA", "North America"
	fmt.Printf("this is my country : %s & this is my continent : %s\n", country, continent)

	// multiple declaration different types of the variables:
	var (
		isEmployed bool   = false
		Salary     int    = 60000
		position   string = "developer"
	)

	fmt.Printf("Am i employed : %t, my salary is : %d, my position is : %s\n", isEmployed, Salary, position)

	// zero values or default values:

	var defaultint int
	var defaultString string
	var defaultBool bool
	var defaultFloat float64

	// this (zeroes values or the default valuews) is way with which compiler assigns us that this is int, string, bool or a float.
	// cons : sometimes variables must be there, but the value it has is undefined or does not exist which is not possible due to the default values.

	fmt.Printf("Default values : Integer: %d, String: %s, Boolean: %t, Float: %f\n", defaultint, defaultString, defaultBool, defaultFloat)

	// constant and enums
	const pi = 3.14 // compiler don't give error if we don't use the constant. it matters with the memory allocation.

	const (
		Monday    = 1
		Tuesday   = 2
		Wednesday = 3
	)

	fmt.Printf("Monday: %d - Tuesday %d Wednesday %d\n", Monday, Tuesday, Wednesday)

	// Go does not have the enum so there is a way we can have that using const
	const (
		Jan   int = iota + 1 // 1
		Feb                  //2
		March                //3
		April                // 4
	) // so if the first var is defined using the iota then consecutive values will be incremental with one starting from the current value of the starting variable and iota is ) by default.

	fmt.Printf("Jan : %d, Feb : %d, March: %d, April: %d\n", Jan, Feb, March, April)

	result := Add(3, 4)

	fmt.Println("Result : ", result)

	sum, prod := CalculateSumAndProduct(3,5);
	fmt.Printf("Sum and Product is : %d and %d ", sum , prod)

}

func Add(a, b int) int {
	return a + b
}

func CalculateSumAndProduct(a, b int) (int, int) {
	return a + b, a * b
}
