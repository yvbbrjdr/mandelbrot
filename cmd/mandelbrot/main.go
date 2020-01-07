package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"math"
	"math/big"
	"os"

	"github.com/yvbbrjdr/mandelbrot/pkg/mandelbrot"
)

func processData(a []int) []byte {
	var m int
	for k, v := range a {
		if k == 0 || v > m {
			m = v
		}
	}
	base := math.Log(float64(m))
	ret := make([]byte, len(a))
	for k, v := range a {
		ret[k] = byte(math.Log(float64(v)) * 255 / base)
	}
	return ret
}

func genMandelbrot(xMin, yMin, xSpan *big.Float, width, height, maxIter int,
	julia *mandelbrot.Complex, output string) error {
	data := mandelbrot.MandelbrotPlane(xMin, yMin, xSpan, width, height,
		maxIter, julia)
	pixel := processData(data)
	im := image.NewGray(image.Rect(0, 0, width, height))
	im.Pix = pixel
	f, err := os.Create(output)
	if err != nil {
		return err
	}
	defer f.Close()
	err = png.Encode(f, im)
	if err != nil {
		return err
	}
	return nil
}

func fail(msg string) {
	fmt.Fprintln(os.Stderr, msg)
	os.Exit(1)
}

func main() {
	xMinStr := flag.String("xmin", "-2", "minimum value of x")
	yMinStr := flag.String("ymin", "-2", "minimum value of y")
	xSpanStr := flag.String("xspan", "4", "span of x")
	width := flag.Int("width", 1000, "width of the image")
	height := flag.Int("height", 1000, "height of the image")
	maxIter := flag.Int("maxiter", 512, "maximum number of iterations")
	xJuliaStr := flag.String("xjulia", "", "x coordinate for Julia set")
	yJuliaStr := flag.String("yjulia", "", "y coordinate for Julia set")
	output := flag.String("output", "output.png", "output filename")
	help := flag.Bool("h", false, "print this help message")
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	xMin, ok := big.NewFloat(0).SetString(*xMinStr)
	if !ok {
		fail("invalid argument: xmin")
	}
	yMin, ok := big.NewFloat(0).SetString(*yMinStr)
	if !ok {
		fail("invalid argument: ymin")
	}
	xSpan, ok := big.NewFloat(0).SetString(*xSpanStr)
	if !ok {
		fail("invalid argument: xspan")
	}
	var julia *mandelbrot.Complex
	if *xJuliaStr != "" || *yJuliaStr != "" {
		xJulia, ok := big.NewFloat(0).SetString(*xJuliaStr)
		if !ok {
			fail("invalid argument: xjulia")
		}
		yJulia, ok := big.NewFloat(0).SetString(*yJuliaStr)
		if !ok {
			fail("invalid argument: yjulia")
		}
		julia = mandelbrot.NewComplex(xJulia, yJulia)
	}
	if err := genMandelbrot(xMin, yMin, xSpan, *width, *height, *maxIter, julia,
		*output); err != nil {
		fail(err.Error())
	}
}
