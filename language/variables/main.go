// https://github.com/ardanlabs/gotraining/blob/master/topics/language/variables/README.md

// Declare three variables that are initialized to their zero value and three
// declared with a literal value. Declare variables of type string, int and
// bool. Display the values of those variables.
//
// Declare a new variable of type float32 and initialize the variable by
// converting the literal value of Pi (3.14).
package main

import "fmt"

// main is the entry point for the application.
func main() {

	// Declare variables that are set to their zero value.
	var zero1 string
	var zero2 int
	var zero3 bool

	// Display the value of those variables.
	fmt.Printf("zero1: %v\n", zero1)
	fmt.Printf("zero2: %v\n", zero2)
	fmt.Printf("zero3: %v\n", zero3)

	// Declare variables and initialize.
	// Using the short variable declaration operator.
	lit1 := "string"
	lit2 := 123
	lit3 := true

	// Display the value of those variables.
	fmt.Printf("lit1: %v\n", lit1)
	fmt.Printf("lit2: %v\n", lit2)
	fmt.Printf("lit3: %v\n", lit3)

	// Perform a type conversion.
	var convert float32
	convert = float32(3.14)

	// Display the value of that variable.
	fmt.Printf("convert: %v\n", convert)
}
