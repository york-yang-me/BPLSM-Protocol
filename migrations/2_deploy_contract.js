const BPLSM = artifacts.require("BPLSM");

module.exports = function(deployer) {
  deployer.deploy(BPLSM);
}