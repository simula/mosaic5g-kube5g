package main

import (
	"fmt"
	"strconv"
)

type Salary struct {
	Basic, HRA, TA float64
}

type Employee struct {
	FirstName, LastName, Email string
	Age                        int
	MonthlySalary              []Salary
}

func main() {
	e := Employee{
		FirstName: "Mark",
		LastName:  "Jones",
		Email:     "mark@gmail.com",
		Age:       25,
		MonthlySalary: []Salary{
			Salary{
				Basic: 15000.00,
				HRA:   5000.00,
				TA:    2000.00,
			},
			Salary{
				Basic: 16000.00,
				HRA:   5000.00,
				TA:    2100.00,
			},
			Salary{
				Basic: 17000.00,
				HRA:   5000.00,
				TA:    2200.00,
			},
		},
	}
	fmt.Println(e.FirstName, e.LastName)
	fmt.Println(e.Age)
	fmt.Println(e.Email)
	for i := 0; i < len(e.MonthlySalary); i++ {
		fmt.Println((e.MonthlySalary[i]).Basic)

	}

	value := "123"
	number, err := strconv.ParseUint(value, 10, 32)
	number = number - 1
	lineNumber := strconv.Itoa(int(number - 1))
	fmt.Println(number)
	fmt.Println(lineNumber)
	fmt.Println(err)
}
