package main

import (
	"fmt"
	"math/big"
)

func main() {
	n := 3
	primes := Primes()
	fmt.Println("Simulating ", n, " players...")

	pieces := make([]PedersenPiece, n)
	for j := 0; j < n; j++ {
		pieces[j] = GeneratePedersenPiece(n)
		pieces[j].OtherShares = make(map[int]Pair, n)
		pieces[j].OtherCommitments = make(map[int][]*big.Int, n)
		for k := 0; k < n; k++ {
			pieces[j].OtherCommitments[k] = make([]*big.Int, n)
		}
	}
	//
	//Now, each piece i has n shares
	//each participant j gets a share pieces[i].Shares[j]
	//the commitments are shared with everyone
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j != i {
				pieces[j].OtherShares[i] = Pair{
					X: (&big.Int{}).Set(pieces[i].Shares[j].X),
					Y: (&big.Int{}).Set(pieces[i].Shares[j].Y),
				}
				copy(pieces[j].OtherCommitments[i], pieces[i].Commitments)
			}
		}
	}

	for i := 0; i < n; i++ {
		fmt.Println("Player ", i, ":")
		for j := 0; j < n; j++ {
			if j != i {
				fmt.Print("Verifying player ", j, ": ")
				fmt.Println(
					VerifyCommitment(pieces[i].OtherCommitments[j],
						pieces[i].OtherShares[j], primes.G, primes.P))
			}
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				pieces[i].FinalShare.Add(pieces[i].FinalShare, pieces[i].OtherShares[j].Y)
			}
		}
		pieces[i].FinalShare.Mod(pieces[i].FinalShare, primes.Q)
		fmt.Println("Final share for player ", i, " : ", pieces[i].FinalShare)
	}

	finalPublicKey := big.NewInt(1)
	for i := 0; i < n; i++ {
		finalPublicKey.Mul(finalPublicKey, pieces[i].Commitments[0])
	}
	finalPublicKey.Mod(finalPublicKey, primes.Q)

	fmt.Println("PK: ", finalPublicKey, " (", finalPublicKey.BitLen(), " bits)")

	m := big.NewInt(2528)
	cipher := GenerateCiphertext(m, finalPublicKey)

	sk := big.NewInt(0)
	for i := 0; i < n; i++ {
		sk.Add(sk, pieces[i].Secret)
	}
	sk.Mod(sk, primes.Q)

	fmt.Println("SK: ", sk, " (", sk.BitLen(), " bits)")

	result := DecryptCiphertext(cipher, sk)
	fmt.Println("Result: ", result, " Expected: 2528")
}
