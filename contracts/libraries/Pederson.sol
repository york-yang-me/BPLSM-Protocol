pragma solidity ^0.4.19;

import "./Curve.sol";

library Pederson {
    using Curve for *;
    /**
     * @dev        Verify a pederson commitment by reconstructing commitment using an unsafe method (h should be hidden)
     * @param      commitment  The commitment
     * @param      ms          Message committed to
     * @param      h           Input to HashToPoint
     * @param      seeds       Random value
     * @return     res         Success or failure
     */
    function verifyMultiUnsafe(uint256[2] commitment, uint256[3] ms, uint256[3] seeds, uint256 h) internal returns(bool res)  {
        // Use random point initially to generate 2nd generator H
        Curve.G1Point memory H =  Curve.HashToPoint(h);

        // Generate left point seed * H
        Curve.G1Point memory lf1 = Curve.g1mul(H, seeds[0]);
        Curve.G1Point memory lf2 = Curve.g1mul(H, seeds[1]);
        Curve.G1Point memory lf3 = Curve.g1mul(H, seeds[2]);

        // Generate right point m * G
        Curve.G1Point memory rt1 = Curve.g1mul(Curve.P1(), ms[0]);
        Curve.G1Point memory rt2 = Curve.g1mul(Curve.P1(), ms[1]);
        Curve.G1Point memory rt3 = Curve.g1mul(Curve.P1(), ms[2]);

        // Generate C = m * G + seed * H
        Curve.G1Point memory c = Curve.g1add(Curve.g1add(Curve.g1add(lf1, rt1), Curve.g1add(lf2, rt2)), Curve.g1add(lf3, rt3));

        return (c.X == commitment[0] && c.Y == commitment[1]);
    }
}
