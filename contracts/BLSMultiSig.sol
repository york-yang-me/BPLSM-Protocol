pragma solidity ^0.4.17;

import "contracts/BLS.sol";

contract BLSMultiSig is BLS {
  G1[] pairKeys;
  uint[] weights;
  uint threshold;
  uint state;

  constructor(uint _threshold, uint[] pairKeyX, uint[] pairKeyY, uint[] _weights) public {
    setStateInternal(0, _threshold, pairKeyX, pairKeyY, _weights);
  }

  function updateState(uint numSigners, uint[] newState, bytes signers,
    uint sigX, uint sigY,
    uint pkXi, uint pkXr, uint pkYi, uint pkYr) public returns (bool) {
      require(checkSig(signers, newState, sigX, sigY, pkXi, pkXr, pkYi, pkYr));
      require(newState.length == 3*numSigners + 2);
      require(newState[0] > state);
      uint[] memory pairKeyX = new uint[](numSigners);
      uint[] memory pairKeyY = new uint[](numSigners);
      uint[] memory _weights = new uint[](numSigners);
      for (uint i = 0; i < numSigners; i++) {
        pairKeyX[i] = newState[i*3+2];
        pairKeyY[i] = newState[i*3+3];
        _weights[i] = newState[i*3+4];
      }
      setStateInternal(newState[0], newState[1], pairKeyX, pairKeyY, _weights);
      return true;
  }

  function setStateInternal(uint _state, uint _threshold, uint[] pairKeyX, uint[] pairKeyY, uint[] _weights) internal {
    assert(pairKeyX.length == pairKeyY.length && pairKeyX.length == _weights.length);
    pairKeys.length = pairKeyX.length;
    for (uint i = 0; i < pairKeyX.length; i++) {
      pairKeys[i] = G1(pairKeyX[i], pairKeyY[i]);
    }
    weights = _weights;
    threshold = _threshold;
    state = _state;
  }

  function isQuorum(bytes signers) public view returns (bool){
    uint weight = 0;
    for (uint i = 0; i < weights.length; i++) {
      if (chkBit(signers,i)) {
        weight += weights[i];
      }
    }
    return weight >= threshold;
  }

  function checkAggKey(bytes signers, G2 aggKey) internal returns (bool) {
    return pairingCheck(sumPoints(pairKeys, signers),g2,g1,aggKey);
  }

  function checkSig(bytes signers, uint[] message,
    uint sigX, uint sigY,
    uint pkXi, uint pkXr, uint pkYi, uint pkYr) public returns (bool) {
      G2 memory aggKey = G2(pkXi, pkXr, pkYi, pkYr);
      G1 memory sig = G1(sigX, sigY);
      return isQuorum(signers) &&
            checkAggKey(signers, aggKey) &&
            checkSignature(message, sig, aggKey);
  }
}
