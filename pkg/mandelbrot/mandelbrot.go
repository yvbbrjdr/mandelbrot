package mandelbrot

import (
	"fmt"
	"math/big"
	"os"
	"runtime"
)

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

func Row(xMin, xSpan, y *big.Float, width, maxIter int, julia *Complex) []int {
	fWidth := big.NewFloat(float64(width))

	ret := make([]int, width)
	for j := 0; j < width; j++ {
		x := big.NewFloat(float64(j))
		x.Mul(x, xSpan).Quo(x, fWidth).Add(x, xMin)

		ret[j] = Mandelbrot(NewComplex(x, y), maxIter, julia)
	}

	return ret
}

func Plane(xMin, yMin, xSpan *big.Float, width, height, maxIter int,
	julia *Complex) []int {
	type dataPacket struct {
		offset int
		data   []int
	}

	fWidth := big.NewFloat(float64(width))
	fHeight := big.NewFloat(float64(height))
	ySpan := big.NewFloat(0).Mul(xSpan, fHeight)
	ySpan.Quo(ySpan, fWidth)

	dataChan := make(chan dataPacket, height)
	signalChan := make(chan struct{}, runtime.NumCPU())
	for i := 0; i < height; i++ {
		fmt.Fprintf(os.Stderr, "Computing Row %d/%d\r", i+1, height)

		y := big.NewFloat(float64(height - i - 1))
		y.Mul(y, ySpan).Quo(y, fHeight).Add(y, yMin)

		offset := i * width
		signalChan <- struct{}{}
		go func() {
			dataChan <- dataPacket{
				offset,
				Row(xMin, xSpan, y, width, maxIter, julia),
			}
			<-signalChan
		}()
	}
	fmt.Fprintln(os.Stderr)

	ret := make([]int, width*height)
	for i := 0; i < height; i++ {
		data := <-dataChan
		copy(ret[data.offset:], data.data)
	}

	var m int
	for k, v := range ret {
		if k == 0 || v > m {
			m = v
		}
	}
	m++
	for k, v := range ret {
		if v == -1 {
			ret[k] = m
		}
	}

	return ret
}
