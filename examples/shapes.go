package examples

import (
	"fmt"
	"github.com/bernsblack/fiddleware/util"
	"math"
)

type Shape interface {
	Area() float64
}

type Square struct {
	Width  float64
	Height float64
}

type Circle struct {
	Radius float64
}

func (c *Square) Area() float64 {
	fmt.Println("functionName:", util.CurrentFunction())
	return c.Height * c.Width
}

func (c *Circle) Area() float64 {
	fmt.Println("functionName:", util.CurrentFunction())
	return math.Pi * math.Pi * c.Radius
}
