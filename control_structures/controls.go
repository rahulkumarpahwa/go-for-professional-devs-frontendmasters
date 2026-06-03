package main

import "fmt"

func main() {
	age := 30

	if age >= 18 {
		fmt.Println("You are an adult!")
	} else if age >= 13 {
		fmt.Println("You are a teenager!")
	} else {
		fmt.Println("You are a child!")
	}

	day := "Tuesday"
	switch day {
	case "Monday":
		fmt.Println("Monday!")
	case "Tuesday", "Wednesday", "Thursday":
		fmt.Println("Midweek!")
	case "Friday":
		fmt.Println("TGIF!")
	case "Saturday", "Sunday":
		fmt.Println("Weekends")
		fallthrough // will go to the consecutive case and not break
	default:
		fmt.Println("Enter a valid!")
	}

	//loops:

	for i := 0; i < 5; i++ {
		fmt.Println("this is i: ", i)
	}

	counter := 0;
	for counter < 3{
		fmt.Println("Counter: ", counter)
		counter++;
	}

	// infinite loop
	interations := 0
	for {
		if interations > 3 {
			break;
		}
		interations++
 	} 

}
