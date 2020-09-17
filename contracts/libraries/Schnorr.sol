pragma solidity ^0.4.19;

import "./Curve.sol";

library Schnorr
{
	function VerifyMultiSchnorr(uint256 signs, uint256[2] R, uint256[2] pb_key, uint256 sumCom) constant internal returns (bool) {
		Curve.G1Point memory sG = Curve.g1mul(Curve.G(), signs);
		Curve.G1Point memory sV = Curve.g1add(
			Curve.G1Point(R[0], R[1]),
			Curve.g1mul(
				Curve.G1Point(pb_key[0], pb_key[1]), sumCom
			)
		);

		return sG.X == sV.X && sG.Y == sV.Y;
	}

}