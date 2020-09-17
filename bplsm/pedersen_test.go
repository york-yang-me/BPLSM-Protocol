// Copyright (c) 2018 Ton van de Ven
// update by york-yang
package bplsm

import (
	"github.com/bwesterb/go-ristretto"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

// Should commit to a sum of two values
func TestCommitToSuccess(t *testing.T) {

	var rX, rY, vX, vY ristretto.Scalar
	rX.Rand()
	H := generateH() // Secondary point on the Curve
	five := big.NewInt(5)

	// Transfer amount of 5 tokens
	tC := commitTo(&H, &rX, vX.SetBigInt(five))

	// Alice 10 - 5 = 5
	rY.Rand()
	ten := big.NewInt(10)
	aC1 := commitTo(&H, &rY, vY.SetBigInt(ten))
	assert.NotEqual(t, aC1, tC, "Should not be equal")
	var aC2 ristretto.Point
	aC2.Sub(&aC1, &tC)

	checkAC2 := SubPrivately(&H, &rX, &rY, ten, five)
	assert.True(t, checkAC2.Equals(&aC2), "Should be equal")
}

// Should fail if not using the correct blinding factors
func TestCommitToFails(t *testing.T) {

	var rX, rY, vX, vY ristretto.Scalar
	rX.Rand()
	H := generateH() // Secondary point on the Curve
	five := big.NewInt(5)

	// Transfer amount of 5 tokens
	tC := commitTo(&H, &rX, vX.SetBigInt(five))

	// Alice 10 - 5 = 5
	rY.Rand()
	ten := big.NewInt(10)
	aC1 := commitTo(&H, &rY, vY.SetBigInt(ten))
	assert.NotEqual(t, aC1, tC, "They should not be equal")
	var aC2 ristretto.Point
	aC2.Sub(&aC1, &tC)

	// Create different (and wrong) binding factors
	rX.Rand()
	rY.Rand()
	checkAC2 := SubPrivately(&H, &rX, &rY, ten, five)
	assert.False(t, checkAC2.Equals(&aC2), "Should not be equal")
}
