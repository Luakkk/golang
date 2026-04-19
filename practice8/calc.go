package main

import "errors"

// Add returns the sum of two integers.
func Add(a, b int) int {
	return a + b
}

// Task 1: Subtract returns the difference of two integers.
func Subtract(a, b int) int {
	return a - b
}

// Task 1: Divide divides a by b. Returns error on division by zero.
func Divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}
