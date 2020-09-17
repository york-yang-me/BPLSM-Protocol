// rewritten by york-yang according to KeyFuse Labs' code and Ven's code
package bplsm

import (
	"github.com/bwesterb/go-ristretto"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

/**
unit test
*/
func TestBPLSM(t *testing.T) {
	hash := DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	// Party 1.
	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := PrvKeyFromBytes(p1.Bytes())
	pub1 := prv1.PubKey()
	party1, err := NewBPLSMParty(prv1)
	assert.Nil(t, err)
	defer party1.Close()

	// Party 2.
	p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	prv2 := PrvKeyFromBytes(p2.Bytes())
	pub2 := prv2.PubKey()
	party2, err := NewBPLSMParty(prv2)
	assert.Nil(t, err)
	defer party2.Close()

	// Party 3.
	p3, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c8", 16)
	prv3 := PrvKeyFromBytes(p3.Bytes())
	pub3 := prv3.PubKey()
	party3, err := NewBPLSMParty(prv3)
	assert.Nil(t, err)
	defer party3.Close()

	// Phase 1.
	sharepub1 := party1.Phase1(pub2).Add(pub3)
	sharepub2 := party2.Phase1(pub1).Add(pub3)
	sharepub3 := party3.Phase1(pub1).Add(pub2)
	assert.Equal(t, sharepub1, sharepub2, sharepub3)

	// Phase 2.
	r1 := party1.Phase2(hash)
	r2 := party2.Phase2(hash)
	r3 := party3.Phase2(hash)

	// Phase 3.
	sharer1 := party1.Phase3(r2).Add(party1.curve, r3)
	sharer2 := party2.Phase3(r1).Add(party2.curve, r3)
	sharer3 := party3.Phase3(r1).Add(party3.curve, r2)
	assert.Equal(t, sharer1, sharer2, sharer3)

	// Phase 4.
	s1, err := party1.Phase4(sharepub1, sharer1)
	assert.Nil(t, err)
	s2, err := party2.Phase4(sharepub2, sharer2)
	assert.Nil(t, err)
	s3, err := party3.Phase4(sharepub3, sharer3)
	assert.Nil(t, err)

	// Phase 5.
	fs1, err := party1.Phase5(sharer1, s1, s2, s3)
	assert.Nil(t, err)
	fs2, err := party2.Phase5(sharer2, s1, s2, s3)
	assert.Nil(t, err)
	fs3, err := party3.Phase5(sharer3, s1, s2, s3)
	assert.Nil(t, err)
	assert.Equal(t, fs1, fs2, fs3)
}

func TestBPLSMCommit(t *testing.T) {
	// Party 1.
	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := PrvKeyFromBytes(p1.Bytes())
	party1, err := NewBPLSMParty(prv1)
	assert.Nil(t, err)
	defer party1.Close()

	// Party 2.
	p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	prv2 := PrvKeyFromBytes(p2.Bytes())
	party2, err := NewBPLSMParty(prv2)
	assert.Nil(t, err)
	defer party2.Close()

	// Party 3.
	p3, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c8", 16)
	prv3 := PrvKeyFromBytes(p3.Bytes())
	party3, err := NewBPLSMParty(prv3)
	assert.Nil(t, err)
	defer party3.Close()

	var seed1, seed2, seed3, v1, v2, v3 ristretto.Scalar
	seed1.Rand()
	G := BPLSGenerateG() // Secondary point on the Curve
	m1 := big.NewInt(1)
	m2 := big.NewInt(2)
	m3 := big.NewInt(3)

	C1, _ := party1.BPLSMCommit(&G, &seed1, v1.SetBigInt(m1))
	C2, _ := party2.BPLSMCommit(&G, &seed2, v2.SetBigInt(m2))
	C3, _ := party3.BPLSMCommit(&G, &seed3, v3.SetBigInt(m3))

	commitments, err := party1.BPLSMSumCommit(C1, C2, C3)
	commitments2, err := party2.BPLSMSumCommit(C1, C2, C3)
	commitments3, err := party3.BPLSMSumCommit(C1, C2, C3)
	assert.Nil(t, err)
	assert.Equal(t, commitments, commitments2, commitments3)
}

/**
stress test
*/
func BenchmarkBPLSMPedersenG(b *testing.B) {
	// Party 1.
	p, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv := PrvKeyFromBytes(p.Bytes())
	party, _ := NewBPLSMParty(prv)
	defer party.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		BPLSGenerateG() // Secondary point on the Curve
	}
}

func BenchmarkBPLSMKeyGen(b *testing.B) {
	// Party 1.
	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := PrvKeyFromBytes(p1.Bytes())
	pub1 := prv1.PubKey()
	party1, _ := NewBPLSMParty(prv1)
	defer party1.Close()

	// Party 2.
	p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	prv2 := PrvKeyFromBytes(p2.Bytes())
	pub2 := prv2.PubKey()
	party2, _ := NewBPLSMParty(prv2)
	defer party2.Close()

	// Party 3.
	p3, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c8", 16)
	prv3 := PrvKeyFromBytes(p3.Bytes())
	party3, _ := NewBPLSMParty(prv3)
	defer party3.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		party3.Phase1(pub1).Add(pub2)
	}
}

func BenchmarkBPLSMSumCom(b *testing.B) {
	// Party 1.
	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := PrvKeyFromBytes(p1.Bytes())
	party1, _ := NewBPLSMParty(prv1)
	defer party1.Close()

	// Party 2.
	p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	prv2 := PrvKeyFromBytes(p2.Bytes())
	party2, _ := NewBPLSMParty(prv2)
	defer party2.Close()

	// Party 3.
	p3, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c8", 16)
	prv3 := PrvKeyFromBytes(p3.Bytes())
	party3, _ := NewBPLSMParty(prv3)
	defer party3.Close()

	var seed1, seed2, seed3, v1, v2, v3 ristretto.Scalar
	seed1.Rand()
	G := BPLSGenerateG() // Secondary point on the Curve
	m1 := big.NewInt(1)
	m2 := big.NewInt(2)
	m3 := big.NewInt(3)

	C1, _ := party1.BPLSMCommit(&G, &seed1, v1.SetBigInt(m1))
	C2, _ := party2.BPLSMCommit(&G, &seed2, v2.SetBigInt(m2))
	C3, _ := party3.BPLSMCommit(&G, &seed3, v3.SetBigInt(m3))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// assume that party1 is the last one to sign and try to aggregate signatures
		if _, err := party1.BPLSMSumCommit(C1, C2, C3); err != nil {
			panic(err)
		}
	}
}

func BenchmarkBPLSMSignature(b *testing.B) {
	hash := DoubleSha256([]byte{0x01, 0x02, 0x03, 0x04})

	// Party 1.
	p1, _ := new(big.Int).SetString("15bafcb56279dbfd985d4d17cdaf9bbfc6701b628f9fb00d6d1e0d2cb503ede3", 16)
	prv1 := PrvKeyFromBytes(p1.Bytes())
	pub1 := prv1.PubKey()
	party1, _ := NewBPLSMParty(prv1)
	defer party1.Close()

	// Party 2.
	p2, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c9", 16)
	prv2 := PrvKeyFromBytes(p2.Bytes())
	pub2 := prv2.PubKey()
	party2, _ := NewBPLSMParty(prv2)
	defer party2.Close()

	// Party 3.
	p3, _ := new(big.Int).SetString("76818c328b8aa1e8f17bd599016fef8134b7d5ec315e0b6373953da7e8b5c0c8", 16)
	prv3 := PrvKeyFromBytes(p3.Bytes())
	pub3 := prv3.PubKey()
	party3, _ := NewBPLSMParty(prv3)
	defer party3.Close()

	// Phase 1.
	sharepub1 := party1.Phase1(pub2).Add(pub3)
	sharepub2 := party2.Phase1(pub1).Add(pub3)
	sharepub3 := party3.Phase1(pub1).Add(pub2)

	// Phase 2.
	r1 := party1.Phase2(hash)
	r2 := party2.Phase2(hash)
	r3 := party3.Phase2(hash)

	// Phase 3.
	sharer1 := party1.Phase3(r2).Add(party1.curve, r3)
	sharer2 := party2.Phase3(r1).Add(party2.curve, r3)
	sharer3 := party3.Phase3(r1).Add(party3.curve, r2)

	// Phase 4.
	s1, _ := party1.Phase4(sharepub1, sharer1)
	s2, _ := party2.Phase4(sharepub2, sharer2)
	s3, _ := party3.Phase4(sharepub3, sharer3)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Phase 5.
		if _, err := party1.Phase5(sharer1, s1, s2, s3); err != nil {
			panic(err)
		}
	}
}
