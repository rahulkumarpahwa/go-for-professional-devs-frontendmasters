package main

import "fmt"

// Structs are collections of fields. They can be used to group data together to form records.
// it is not only held in memory but also can be used to define methods on them. Structs are useful for grouping data together to form records.
// it also act as the type which can be used to define methods on them. It is a way to create complex data types that group together related data and behavior. Structs can be used to model real-world entities, such as people, cars, or products, and can be used to define methods that operate on the data they contain.
type Person struct {
	Name string
	Age  int
}

// nested struct
type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

// Structs can also be nested, meaning that a struct can contain another struct as a field. This allows for more complex data structures to be created. For example, we can create a struct called "Employee" that contains a "Person" struct and an "Address" struct as fields.

type contactInfo struct {
	Email   string
	Phone   string
	Address Address // nested struct
}

func main() {
	person1 := Person{Name: "John Doe" + " " + "Smith", Age: 30}
	fmt.Println("Person 1:", person1)

	person2 := Person{Name: "Jane Doe", Age: 25}
	fmt.Printf("Person 2: %+v\n", person2) // + sign is used to print the field names along with their values

	// Accessing fields of a struct
	fmt.Println("Person 1 Name:", person1.Name)
	fmt.Println("Person 1 Age:", person1.Age)

	// Modifying fields of a struct
	person1.Age = 31
	fmt.Println("Person 1 after birthday:", person1)

	// anonymous struct
	// Anonymous structs are useful when you need to create a struct type that is only used in one place and doesn't need to be reused elsewhere in your code. They can help you avoid cluttering your code with unnecessary struct definitions and can make your code more concise and easier to read.
	// they are used in testing when you want to create a struct type that is only used in one test case and doesn't need to be reused elsewhere in your code. They can help you avoid cluttering your test code with unnecessary struct definitions and can make your tests more concise and easier to read.
	anonymousPerson := struct {
		Name string
		Age  int
	}{
		Name: "Anonymous",
		Age:  40,
	}
	fmt.Println("Anonymous Person:", anonymousPerson)

	// usage of nested struct
	contact := contactInfo{
		Email: "apple@apple.com",
		Phone: "123-456-7890",
		Address: Address{
			Street: "123 Main St",
			City:   "Anytown",
			State:  "CA",
			Zip:    "12345",
		},
	}
	fmt.Println("Contact Info:", contact)


	// Using pointers to structs
	fmt.Println("Person 1 before modification:", person1)
	personPointer := &person1 // creating a pointer to the person1 struct
	ModifyPerson(personPointer)

	x := 20	
	ptr := &x // creating a pointer to the variable x
	fmt.Println("Value of x before modification:", x)
	fmt.Println("Pointer to x:", ptr)
	*ptr = 30 // modifying the value of x using the pointer
	fmt.Println("Value of x after modification:", x)

	// structs allow us to define methods on them, which can be used to modify the fields of the struct. This is done by defining a method with a receiver of the struct type. The receiver is a special parameter that allows the method to access and modify the fields of the struct. In this example, we define a method called "modifyPersonName" that takes a string parameter and modifies the Name field of the Person struct.
	// Using methods with structs
	person1.modifyPersonName("John Smith")
	fmt.Println("Person 1 after name modification:", person1)

}

func ModifyPerson(p *Person) {
	// Modifying fields of a struct using a pointer
	p.Age = 35
	fmt.Println("Person after modification:", *p)
}

func (p *Person) modifyPersonName(name string) {
	p.Name = name
	fmt.Println("Person name after modification:", p.Name)
}
