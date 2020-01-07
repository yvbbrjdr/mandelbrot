package mandelbrot

import (
	"fmt"
	"math/big"
	"os"
)

type dataPacket struct {
	index int
	data  int
}

func Mandelbrot(p *Complex, maxIter int, julia *Complex) int {
	var x, c *Complex
	if julia == nil {
		x = NewComplex(big.NewFloat(0), big.NewFloat(0))
		c = p
	} else {
		x = p.Copy()
		c = julia
	}
	for i := 1; i <= maxIter; i++ {
		if x.SqrAndInc(c).BlowUp() {
			return i
		}
	}
	return -1
}

func MandelbrotPlane(xMin, yMin, xSpan *big.Float, width, height, maxIter int,
	julia *Complex) []int {
	fWidth := big.NewFloat(float64(width))
	fHeight := big.NewFloat(float64(height))
	ySpan := big.NewFloat(0).Mul(xSpan, fHeight)
	ySpan.Quo(ySpan, fWidth)
	a := make([]int, width*height)
	for i := 0; i < height; i++ {
		fmt.Fprintf(os.Stderr, "Computing Row %d\r", i)
		y := big.NewFloat(float64(height - i - 1))
		y.Mul(y, ySpan).Quo(y, fHeight).Add(y, yMin)
		c := make(chan dataPacket, width)
		for j := 0; j < width; j++ {
			x := big.NewFloat(float64(j))
			x.Mul(x, xSpan).Quo(x, fWidth).Add(x, xMin)
			index := i*width + j
			go func() {
				c <- dataPacket{
					index,
					Mandelbrot(NewComplex(x, y), maxIter, julia),
				}
			}()
		}
		for j := 0; j < width; j++ {
			packet := <-c
			a[packet.index] = packet.data
		}
	}
	fmt.Fprintln(os.Stderr)
	var m int
	for k, v := range a {
		if k == 0 || v > m {
			m = v
		}
	}
	for k, v := range a {
		if v == -1 {
			a[k] = m + 1
		}
	}
	return a
}
