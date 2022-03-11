package main

import (
	"crypto/rand"
	"math/big"
)

// GermainPrime returns a prime q such that 2q + 1 is also prime
func GermainPrime(bits int) *big.Int {
	for {
		q, _ := rand.Prime(rand.Reader, bits)
		p := (&big.Int{}).Mul(q, big.NewInt(2))
		p.Add(p, big.NewInt(1))
		if p.ProbablyPrime(40) {
			return q
		}
	}
}

// FeldmanPrimes returns three primes q, p, g where
// q is the prime used for polynomial interpolation
// p is 2q + 1
// g is the generator used for commitments
func FeldmanPrimes(bits int) PrimeTriple {
	q := GermainPrime(bits)
	p := (&big.Int{}).Set(q)
	p.Mul(p, big.NewInt(2))
	p.Add(p, big.NewInt(1))

	b, _ := rand.Int(rand.Reader, p)
	for b.BitLen() == 0 {
		b, _ = rand.Int(rand.Reader, p)
	}

	g := (&big.Int{}).Mul(b, b)
	g.Mod(g, p)

	primes := PrimeTriple{
		P: p,
		Q: q,
		G: g,
	}
	return primes
}

// GenerateCommitments generates commitments for feldman
func GenerateCommitments(poly Polynomial, g *big.Int, p *big.Int) []*big.Int {
	commitments := make([]*big.Int, poly.Degree)

	for i := 0; i < poly.Degree; i++ {
		comm := (&big.Int{}).Exp(g, poly.Coefficients[i], p)
		commitments[i] = comm
	}

	return commitments
}

// VerifyCommitment basically does what it says it does
func VerifyCommitment(commitments []*big.Int, share Pair, g *big.Int, p *big.Int) bool {
	gv := (&big.Int{}).Exp(g, share.Y, p)
	prod := (&big.Int{}).Set(commitments[0])
	for j := 1; j < len(commitments); j++ {
		exp := (&big.Int{}).Exp(share.X, big.NewInt(int64(j)), p)
		exp.Exp(commitments[j], exp, p)
		prod.Mul(prod, exp)
	}
	prod.Mod(prod, p)
	return gv.Cmp(prod) == 0
}
