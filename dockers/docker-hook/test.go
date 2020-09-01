package main

import "fmt"

type Color int

const (
	red Color = iota
	yellow
	green
	blue
)

func (c *Color) Name() string {
	switch *c {
	case red:
		return "red"
	case yellow:
		return "yellow"
	case green:
		return "green"
	case blue:
		return "blue"
	}
	return ""
}

func main() {
	col := red
	fmt.Println("col = ", col.Name())
}
