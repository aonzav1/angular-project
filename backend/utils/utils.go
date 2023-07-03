package utils

import (
	"fmt"
)

// Function for handling errors
func CheckErr(err error) {
	if err != nil {
		fmt.Println("========= ERROR CAUGHT ========")
		panic(err)
	}
}

// Function for handling messages
func PrintMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}
