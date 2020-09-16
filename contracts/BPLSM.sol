pragma solidity ^0.4.19;
import "./libraries/Pederson.sol";
import "./libraries/Schnorr.sol";

contract BPLSM {

    using Pederson for *;
    using Schnorr for *;

    function verifyBPLSMSumCommit(uint256[2] Com, uint256[3] declares, uint256[3] seeds, uint256 H) public view returns (bool) {
        return Pederson.verifyMultiUnsafe(Com, declares, seeds, H);
    }

    function verifyBPLSMSchnorr(uint256 signs, uint256[2] R, uint256[2] pb_key, uint256 sumCom) public view returns (bool) {
        return Schnorr.VerifyMultiSchnorr(signs, R, pb_key, sumCom);
    }

}
