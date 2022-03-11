package main

import (
	"math/big"
)

// GenerateShares sets up the split for all participants t of n
func GenerateShares(s *big.Int, primes PrimeTriple, t int, n int) ([]Pair, Polynomial) {
	poly := GenerateRandomPolynomial(t, s, primes)
	return poly.GeneratePoints(big.NewInt(1), big.NewInt(int64(n+1)), n), poly
}

// CombineSecret takes an array of shares and gets the secret
func CombineSecret(points []Pair, p *big.Int) *big.Int {
	return LagrangeInterpolation(points, p)
}
