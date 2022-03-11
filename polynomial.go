package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Pair is an X,Y pair on some field with arbitrary size integers
type Pair struct {
	X, Y *big.Int
}

// Polynomial is a coefficient representation of an arbitrary degree over the modulo p field
type Polynomial struct {
	Coefficients []*big.Int
	Degree       int
	Primes       PrimeTriple
}

// GenerateRandomPolynomial a random polynomial with coefficients between 0 and maxRand mod p
func GenerateRandomPolynomial(degree int, intercept *big.Int, primes PrimeTriple) Polynomial {
	coeffs := make([]*big.Int, degree)
	coeffs[0] = intercept
	for i := 1; i < degree; i++ {
		r, _ := rand.Int(rand.Reader, primes.Q)
		coeffs[i] = r
	}

	return Polynomial{
		Coefficients: coeffs,
		Degree:       degree,
		Primes: PrimeTriple{
			P: (&big.Int{}).Set(primes.P),
			Q: (&big.Int{}).Set(primes.Q),
			G: (&big.Int{}).Set(primes.G),
		},
	}
}

// EvaluatePolynomial evaluates a polynomial at some point
func (poly *Polynomial) EvaluatePolynomial(x *big.Int) *big.Int {
	res := &big.Int{}
	res.Set(poly.Coefficients[0])

	for i := 1; i < poly.Degree; i++ {
		curr := (&big.Int{}).Set(poly.Coefficients[i])
		xpow := (&big.Int{}).Set(x)
		for j := 1; j < i; j++ {
			xpow.Mul(xpow, x)
		}
		curr.Mul(curr, xpow)
		res.Add(res, curr)
	}

	res.Mod(res, poly.Primes.Q)
	return res
}

// String turns a polynomial into a string representation
func (poly *Polynomial) String() string {
	str := ""
	str += fmt.Sprintf("y = %s + %sx", poly.Coefficients[0], poly.Coefficients[1])
	for i := 2; i < poly.Degree; i++ {
		str += fmt.Sprintf(" + %sx^%d", poly.Coefficients[i], i)
	}

	return str
}

// GeneratePoints takes a polynomial and generates n random points from a to b
func (poly *Polynomial) GeneratePoints(a, b *big.Int, n int) []Pair {
	pairs := make([]Pair, n)
	delta := (&big.Int{}).Sub(b, a)
	delta.Div(delta, big.NewInt(int64(n)))
	x := (&big.Int{}).Set(a)

	for i := 0; i < n; i++ {
		newX := (&big.Int{}).Set(x)
		pairs[i] = Pair{
			X: newX,
			Y: poly.EvaluatePolynomial(newX),
		}
		x.Add(x, delta)
	}

	return pairs
}

func LagrangeInterpolation(points []Pair, p *big.Int) *big.Int {
	res := big.NewInt(0)
	for j := 0; j < len(points); j++ {
		num := big.NewInt(1)
		denom := big.NewInt(1)
		for m := 0; m < len(points); m++ {
			if m != j {
				num.Mul(num, points[m].X)
				curr_denom := &big.Int{}
				curr_denom.Sub(points[m].X, points[j].X)
				denom.Mul(denom, curr_denom)
			}
		}
		basis := big.NewInt(1)
		basis.Mul(basis, points[j].Y)
		basis.Mul(basis, num)
		basis.Div(basis, denom)

		res.Add(res, basis)
	}

	res.Mod(res, p)
	return res
}
