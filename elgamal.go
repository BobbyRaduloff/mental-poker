package main

import (
	"crypto/rand"
	"math/big"
)

// Primes taken from: https://www.rfc-editor.org/rfc/rfc5114#section-2.2

// PHex is P for ElGamal
const PHex = "B10B8F96A080E01DDE92DE5EAE5D54EC52C99FBCFB06A3C69A6A9DCA52D23B616073E28675A23D189838EF1E2EE652C013ECB4AEA906112324975C3CD49B83BFACCBDD7D90C4BD7098488E9C219A73724EFFD6FAE5644738FAA31A4FF55BCCC0A151AF5F0DC8B4BD45BF37DF365C1A65E68CFDA76D4DA708DF1FB2BC2E4A4371"

// GHex is G for ElGamal
const GHex = "A4D1CBD5C3FD34126765A442EFB99905F8104DD258AC507FD6406CFF14266D31266FEA1E5C41564B777E690F5504F213160217B4B01B886A5E91547F9E2749F4D7FBD7D3B9A92EE1909D0D2263F80A76A6A24C087A091F531DBF0A0169B6A28AD662A4D18E73AFA32D779D5918D08BC8858F4DCEF97C2A24855E6EEB22B3B2E5"

// QHex is Q for ElGamal
const QHex = "F518AA8781A8DF278ABA4E7D64B7CB9D49462353"

// PublicKey is the ElGamal Public Key g^x
type PublicKey *big.Int

// PrivateKey is the ElGamal Private Key x
type PrivateKey *big.Int

// CipherText is the encrypted message
type CipherText struct {
	A *big.Int
	B *big.Int
}

type PrimeTriple struct {
	P *big.Int
	Q *big.Int
	G *big.Int
}

// Primes gives you the ElGamal primes
func Primes() PrimeTriple {
	q := &big.Int{}
	q.SetString(QHex, 16)
	p := &big.Int{}
	p.SetString(PHex, 16)
	g := &big.Int{}
	g.SetString(GHex, 16)

	return PrimeTriple{
		P: p,
		Q: q,
		G: g,
	}
}

// GenerateKeys generates a keypair
func GenerateKeys() (PrivateKey, PublicKey) {
	primes := Primes()

	sk, _ := rand.Int(rand.Reader, primes.Q)
	pk := (&big.Int{}).Exp(primes.G, sk, primes.P)

	return sk, pk
}

// GenerateCiphertext encrypts m using pk
func GenerateCiphertext(m *big.Int, pk PublicKey) CipherText {
	primes := Primes()

	r, _ := rand.Int(rand.Reader, primes.Q)

	a := &big.Int{}
	a.Exp(primes.G, r, primes.P)

	b := &big.Int{}
	b.Exp(pk, r, primes.P)
	b.Mul(b, m)

	cipher := CipherText{
		A: a,
		B: b,
	}

	return cipher
}

// DecryptCiphertext decrypts c with pk
func DecryptCiphertext(c CipherText, pk PrivateKey) *big.Int {
	primes := Primes()

	m := &big.Int{}
	d := &big.Int{}
	d.Exp(c.A, pk, primes.P)
	m.Div(c.B, d)

	return m
}
