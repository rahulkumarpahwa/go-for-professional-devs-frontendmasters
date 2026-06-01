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
	// cons : sometimes variables must be there, but the value it has undefined or does not exist which is not possible due to the default values.

	fmt.Printf("Default values : Integer: %d, String: %s, Boolean: %t, Float: %f\n", defaultint, defaultString, defaultBool, defaultFloat)
}
