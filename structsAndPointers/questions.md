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

6. What is the difference between passing a struct by value versus passing a struct by pointer in Go?
6. When passing a struct by value, a copy of the struct is created, and modifications are only local to the function's scope. When passing a struct by pointer, the function can modify the original struct directly, and changes persist outside the function's scope.

7. How do you create a pointer to a variable in Go?
7. To create a pointer to a variable, use the ampersand (&) symbol before the variable name, which returns the memory address of the variable. To dereference and modify the value, use the asterisk (*) symbol before the pointer variable.

8. What is a method receiver in Go?
8. A method receiver is a way to define methods on a specific struct type. It allows you to associate a function with a struct by declaring the receiver type in parentheses before the function name, enabling method calls directly on struct instances.

9. What naming convention determines a method or field's visibility in Go?
9. Methods and fields starting with a capital letter are exported and can be accessed from other packages. Methods and fields starting with a lowercase letter are only accessible within the same package.

10. When should you consider using pointers instead of passing values in Go?
10. Consider using pointers when you want to modify the original value across multiple functions, avoid copying large structs, or when you need to persist changes outside a function's scope. A good rule of thumb is to use pointers when passing structs through multiple functions repeatedly.