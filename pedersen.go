package main

import (
	"crypto/rand"
	"math/big"
)

// PedersenPiece is a single user's piece of the Pedersen "puzzle"
// (i.e. one instance of Feldman)
type PedersenPiece struct {
	Secret           *big.Int
	Polynomial       Polynomial
	Shares           []Pair
	Commitments      []*big.Int
	N                int
	OtherShares      map[int]Pair
	OtherCommitments map[int][]*big.Int
	FinalShare       *big.Int
}

// GeneratePedersenPiece creates the instance of Feldman for one player
func GeneratePedersenPiece(n int) PedersenPiece {
	primes := Primes()

	secret, _ := rand.Int(rand.Reader, primes.Q)
	shares, poly := GenerateShares(secret, primes, n, n)
	commitments := GenerateCommitments(poly, primes.G, primes.P)

	piece := PedersenPiece{
		Secret:           secret,
		Polynomial:       poly,
		Shares:           shares,
		Commitments:      commitments,
		N:                n,
		OtherShares:      nil,
		OtherCommitments: nil,
		FinalShare:       big.NewInt(0),
	}

	return piece
}
