// Copyright (c) 2018 Ton van de Ven
// update by york-yang
package BPLSM_protocol

import (
	"github.com/bwesterb/go-ristretto"
	"math/big"
)

// The prime order of the base point is 2^252 + 27742317777372353535851937790883648493.
var n25519, _ = new(big.Int).SetString("7237005577332262213973186563042994240857116359379907606001950938285454250989", 10)

// Commit to a value x
// H - Random secondary point on the curve
// r - Private key used as blinding factor
// x - The value (number of tokens)
func commitTo(H *ristretto.Point, r, x *ristretto.Scalar) ristretto.Point {
	//ec.g.mul(r).add(H.mul(x));
	var result, rPoint, transferPoint ristretto.Point
	rPoint.ScalarMultBase(r)
	transferPoint.ScalarMult(H, x)
	result.Add(&rPoint, &transferPoint)
	return result
}

// Generate a random point on the curve
func generateH() ristretto.Point {
	var random ristretto.Scalar
	var H ristretto.Point
	random.Rand()
	H.ScalarMultBase(&random)
	return H
}

// Subtract two commitments using homomorphic encryption
func Sub(cX, cY *ristretto.Point) ristretto.Point {
	var subPoint ristretto.Point
	subPoint.Sub(cX, cY)
	return subPoint
}

// Subtract two known values with blinding factors
//   and compute the committed value
//   add rX - rY (blinding factor private keys)
//   add vX - vY (hidden values)
func SubPrivately(H *ristretto.Point, rX, rY *ristretto.Scalar, vX, vY *big.Int) ristretto.Point {
	var rDif ristretto.Scalar
	var vDif big.Int
	rDif.Sub(rY, rX)
	vDif.Sub(vX, vY)
	vDif.Mod(&vDif, n25519)

	var vScalar ristretto.Scalar
	var rPoint ristretto.Point
	vScalar.SetBigInt(&vDif)

	rPoint.ScalarMultBase(&rDif)
	var vPoint, result ristretto.Point
	vPoint.ScalarMult(H, &vScalar)
	result.Add(&rPoint, &vPoint)
	return result
}
