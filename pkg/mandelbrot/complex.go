package mandelbrot

import "math/big"

type Complex struct {
	a *big.Float
	b *big.Float
}

func NewComplex(a, b *big.Float) *Complex {
	return &Complex{a, b}
}

func (x *Complex) Copy() *Complex {
	return &Complex{
		big.NewFloat(0).Copy(x.a),
		big.NewFloat(0).Copy(x.b),
	}
}

func (x *Complex) SqrAndInc(c *Complex) *Complex {
	aa := big.NewFloat(0).Mul(x.a, x.a)
	bb := big.NewFloat(0).Mul(x.b, x.b)
	ab := big.NewFloat(0).Mul(x.a, x.b)
	x.a.Sub(aa, bb).Add(x.a, c.a)
	x.b.Add(ab, ab).Add(x.b, c.b)
	return x
}

func (x *Complex) BlowUp() bool {
	a, _ := x.a.Float32()
	b, _ := x.b.Float32()
	return a*a+b*b >= 4
}
