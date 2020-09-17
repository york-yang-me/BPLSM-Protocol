pragma solidity ^0.4.17;

contract BLS {
  struct G1 {
    uint x;
    uint y;
  }
  G1 g1 = G1(1,2);

  struct G2 {
    uint xi;
    uint xr;
    uint yi;
    uint yr;
  }
  G2 g2 = G2(
    11559732032986387107991004021392285783925812861821192530917403151452391805634,
    10857046999023057135944570762232829481370756359578518086990519993285655852781,
    4082367875863433681332203403145435568316851327593401208105741076214120093531,
    8495653923123431417604973247489272438418190587263600148770280649306958101930
  );

  function modPow(uint256 base, uint256 exponent, uint256 modulus) internal returns (uint256) {
    uint256[6] memory input = [32,32,32,base,exponent,modulus];
    uint256[1] memory result;
    assembly {
      if iszero(call(not(0), 0x05, 0, input, 0xc0, result, 0x20)) {
        revert(0, 0)
      }
    }
    return result[0];
  }

  function addPoints(G1 a, G1 b) internal returns (G1) {
    uint256[4] memory input = [a.x, a.y, b.x, b.y];
    uint[2] memory result;
    assembly {
      if iszero(call(not(0), 0x06, 0, input, 0x80, result, 0x40)) {
        revert(0, 0)
      }
    }
    return G1(result[0], result[1]);
  }

  function chkBit(bytes b, uint x) public pure returns (bool) {
    return uint(b[x/8])&(uint(1)<<(x%8)) != 0;
  }

  function sumPoints(G1[] points, bytes indices) internal returns (G1) {
    G1 memory acc = G1(0,0);
    for (uint i = 0; i < points.length; i++) {
      if (chkBit(indices, i)) {
        acc = addPoints(acc, points[i]);
      }
    }
    return G1(acc.x, acc.y);
  }

  function scalarMultiply(G1 point, uint256 scalar) internal returns(G1) {
    uint256[3] memory input = [point.x, point.y, scalar];
    uint[2] memory result;
    assembly {
      if iszero(call(not(0), 0x07, 0, input, 0x60, result, 0x40)) {
        revert(0, 0)
      }
    }
    return G1(result[0], result[1]);
  }

  function pairingCheck(G1 a, G2 x, G1 b, G2 y) internal returns (bool) {
    //returns e(a,x) == e(b,y)
    uint256[12] memory input = [
      a.x, a.y, x.xi, x.xr, x.yi, x.yr, b.x, prime - b.y, y.xi, y.xr, y.yi, y.yr
    ];
    uint[1] memory result;
    assembly {
      if iszero(call(not(0), 0x08, 0, input, 0x180, result, 0x20)) {
          revert(0, 0)
      }
    }
    return result[0]==1;
  }


  uint256 prime = 21888242871839275222246405745257275088696311157297823662689037894645226208583;
  uint256 pminus = 21888242871839275222246405745257275088696311157297823662689037894645226208582;
  uint256 pplus = 21888242871839275222246405745257275088696311157297823662689037894645226208584;

  function hashToG1(uint[] b) internal returns (G1) {
    uint x = 0;
    while (true) {
      uint256 hx = uint256(keccak256(b,byte(x)))%prime;
      uint256 px = (modPow(hx,3,prime) + 3);
      if (modPow(px, pminus/2, prime) == 1) {
        uint256 py = modPow(px, pplus/4, prime);
        if (uint(keccak256(b,byte(255)))%2 == 0)
          return G1(hx,py);
        else
          return G1(hx,prime-py);
      } else {
        x++;
      }
    }
  }

  function checkSignature(uint[] message, G1 sig, G2 aggKey) internal returns (bool) {
    return pairingCheck(sig, g2, hashToG1(message), aggKey);
  }
}
