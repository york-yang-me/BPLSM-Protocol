// rewritten by york-yang according to KeyFuse Labs' code and Ven's code
package bplsm

import (
	"github.com/bwesterb/go-ristretto"
	"math/big"

	"crypto/elliptic"

	"github.com/keyfuse/tokucore/xcrypto/schnorr"
	"github.com/keyfuse/tokucore/xcrypto/secp256k1"
)

// BPLSMParty -- bplsm struct.
type BPLSMParty struct {
	k0    *big.Int
	N     *big.Int
	prv   *PrvKey
	pub   *PubKey
	hash  []byte
	curve elliptic.Curve
	r     *secp256k1.Scalar
}

// NewBPLSMParty -- creates new BPLSMParty.
func NewBPLSMParty(prv *PrvKey) (*BPLSMParty, error) {
	pub := prv.PubKey()
	curve := pub.Curve
	N := curve.Params().N
	return &BPLSMParty{
		N:     N,
		prv:   prv,
		pub:   pub,
		curve: curve,
	}, nil
}

/**
  BPLSM-Pedersen Part
*/
// Commit to message m
// G    - Random secondary point on the curve
// seed - Private key used as blinding factor
// m    - The message
func (party *BPLSMParty) BPLSMCommit(G *ristretto.Point, seed, m *ristretto.Scalar) ([]byte, error) {
	//ec.h.mul(seed).add(G.mul(m));
	var result, sPoint, transferPoint ristretto.Point
	sPoint.ScalarMultBase(seed)
	transferPoint.ScalarMult(G, m)
	// c = m*G + seed*H
	result.Add(&sPoint, &transferPoint)
	return result.Bytes(), nil
}

// Generate a random point G on the curve
func BPLSGenerateG() ristretto.Point {
	var random ristretto.Scalar
	var G ristretto.Point
	random.Rand()
	G.ScalarMultBase(&random)
	return G
}

// return the final commitment
func (party *BPLSMParty) BPLSMSumCommit(coms ...[]byte) ([]byte, error) {
	N := party.N
	aggs := new(big.Int)
	commitments := make([]byte, 64)

	for _, com := range coms {
		c := new(big.Int).SetBytes(com[:])
		aggs.Add(aggs, c)
	}
	aggs = aggs.Mod(aggs, N)

	copy(commitments[:], schnorr.IntToByte(aggs))

	return commitments, nil
}

/**
  BPLSM-Schnorr part
*/
// Phase1 -- used to generate final pubKey of parties.
// Return the shared PubKey.
func (party *BPLSMParty) Phase1(pub2 *PubKey) *PubKey {
	pub := party.pub
	return pub.Add(pub2)
}

// Phase2 -- used to generate k, kInv, scalarR.
// Return the party scalar R.
func (party *BPLSMParty) Phase2(hash []byte) *secp256k1.Scalar {
	N := party.N
	prv := party.prv
	pub := prv.PubKey()
	curve := pub.Curve
	d := schnorr.IntToByte(prv.D)

	party.hash = hash
	// Scalar R.
	// k' = int(hash(bytes(d) || Com)) mod n
	k0, err := schnorr.GetK0(hash, d, N)
	if err != nil {
		return nil
	}
	party.k0 = k0

	rx, ry := curve.ScalarBaseMult(k0.Bytes())
	party.r = secp256k1.NewScalar(rx, ry)
	return party.r
}

// Phase3 -- return shared scalar R.
func (party *BPLSMParty) Phase3(r2 *secp256k1.Scalar) *secp256k1.Scalar {
	curve := party.curve
	scalarR := party.r

	shareScalarR := secp256k1.NewScalar(scalarR.X, scalarR.Y)
	return shareScalarR.Add(curve, r2)
}

// Phase4 -- return the signature of this party.
func (party *BPLSMParty) Phase4(sharePub *PubKey, shareR *secp256k1.Scalar) ([]byte, error) {
	k0 := party.k0
	Com := party.hash
	N := party.N
	prv := party.prv
	pub := sharePub
	curve := party.curve
	scalarR := party.r
	shareScalarR := shareR

	// sumCom = int(hash(bytes(x(R)) || bytes(dG) || Com)) mod n
	sumCom := schnorr.GetE(curve, Com, pub.X, pub.Y, schnorr.IntToByte(scalarR.X))

	// ed
	ed := new(big.Int)
	ed.Mul(sumCom, prv.D)

	// s = k + ed (sig = r + sumCom * pr_key)
	k := schnorr.GetK(curve, shareScalarR.Y, k0)
	s := new(big.Int)
	s.Add(k, ed)
	s.Mod(s, N)

	return schnorr.IntToByte(s), nil
}

// Phase5 -- return the final signature.
func (party *BPLSMParty) Phase5(shareR *secp256k1.Scalar, sigs ...[]byte) ([]byte, error) {
	N := party.N
	R := shareR

	aggs := new(big.Int)
	sigFinal := make([]byte, 64)

	for _, sig := range sigs {
		s := new(big.Int).SetBytes(sig[:])
		aggs.Add(aggs, s)
	}
	aggs = aggs.Mod(aggs, N)

	copy(sigFinal[:32], schnorr.IntToByte(R.X))
	copy(sigFinal[32:], schnorr.IntToByte(aggs))
	return sigFinal, nil
}

// Close -- close the party.
func (party *BPLSMParty) Close() {
	party.prv = nil
	if party.k0 != nil {
		party.k0.SetInt64(0)
	}
}
