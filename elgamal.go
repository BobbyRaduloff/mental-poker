package main

import (
	"crypto/rand"
	"math/big"
)

// Primes taken from: https://www.rfc-editor.org/rfc/rfc5114#section-2.2

// PHex is P for ElGamal
const PHex = "87A8E61DB4B6663CFFBBD19C651959998CEEF608660DD0F25D2CEED4435E3B00E00DF8F1D61957D4FAF7DF4561B2AA3016C3D91134096FAA3BF4296D830E9A7C209E0C6497517ABD5A8A9D306BCF67ED91F9E6725B4758C022E0B1EF4275BF7B6C5BFC11D45F9088B941F54EB1E59BB8BC39A0BF12307F5C4FDB70C581B23F76B63ACAE1CAA6B7902D52526735488A0EF13C6D9A51BFA4AB3AD8347796524D8EF6A167B5A41825D967E144E5140564251CCACB83E6B486F6B3CA3F7971506026C0B857F689962856DED4010ABD0BE621C3A3960A54E710C375F26375D7014103A4B54330C198AF126116D2276E11715F693877FAD7EF09CADB094AE91E1A1597"

// GHex is G for ElGamal
const GHex = "3FB32C9B73134D0B2E77506660EDBD484CA7B18F21EF205407F4793A1A0BA12510DBC15077BE463FFF4FED4AAC0BB555BE3A6C1B0C6B47B1BC3773BF7E8C6F62901228F8C28CBB18A55AE31341000A650196F931C77A57F2DDF463E5E9EC144B777DE62AAAB8A8628AC376D282D6ED3864E67982428EBC831D14348F6F2F9193B5045AF2767164E1DFC967C1FB3F2E55A4BD1BFFE83B9C80D052B985D182EA0ADB2A3B7313D3FE14C8484B1E052588B9B7D2BBD2DF016199ECD06E1557CD0915B3353BBB64E0EC377FD028370DF92B52C7891428CDC67EB6184B523D1DB246C32F63078490F00EF8D647D148D47954515E2327CFEF98C582664B4C0F6CC41659"

// QHex is Q for ElGamal
const QHex = "8CF83642A709A097B447997640129DA299B1A47D1EB3750BA308B0FE64F5FBD3"

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
