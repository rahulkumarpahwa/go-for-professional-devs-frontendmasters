1. What is a struct in Go and how do you define one?
1. A struct is a custom data type that can hold multiple fields of different types. To define a struct, use the 'type' keyword followed by the struct name, the 'struct' keyword, and curly braces containing field names and their types, such as: type Person struct { name string, age int }

2. What is the difference between printing a struct with %v and %+v format specifiers?
2. %v prints the struct values, while %+v prints the struct values with their corresponding field names, providing more explicit output about the struct's contents

3. What is an anonymous struct in Go?
3. An anonymous struct is a struct defined and instantiated in a single place without being given a named type. It can be created directly with fields and values, such as: employee := struct { name string; id int }{ name: "Alice", id: 123 }

4. How can you create a nested struct in Go?
4. A nested struct can be created by defining one struct as a field within another struct. For example, a Contact struct can contain an Address struct as a field, like: type Contact struct { name string; address Address; phone string }

5. How are structs passed in Go by default?
5. In Go, structs are passed by value, which means a complete copy of the struct is created when passed to a function. Modifications to the struct within the function will not affect the original struct unless passed by reference using a pointer