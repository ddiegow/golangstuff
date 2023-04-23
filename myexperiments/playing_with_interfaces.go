package main

import (
	"fmt"
	"math"
)

type shape interface {
	area() float64
	circumf() float64
}

type square struct {
	side float64
}

type circle struct {
	radius float64
}

func (s square) area() float64 {
	return s.side * s.side
}

func (s square) circumf() float64 {
	return s.side * 4
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c circle) circumf() float64 {
	return 2 * math.Pi * c.radius
}

func describe(s shape) {
	fmt.Printf("area: %f\n", s.area())
	fmt.Printf("circumf: %f\n", s.circumf())
}

func main() {
	s := square{5.0}
	s2 := circle{10.0}
	describe(s)
	describe(s2)
}
