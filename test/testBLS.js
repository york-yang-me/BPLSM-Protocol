const BLS = artifacts.require("BLS");
const TestProxy = artifacts.require("BLSTest");

contract("BLS", async (accounts) =>  {
  let bls;
  let blsTest;
  beforeEach(async () => {
    bls = await BLS.new();
    blsTest = await TestProxy.new();
  })
  it("should verify trivial pairing", async () => {
    assert(await blsTest.pairingCheckTrivial.call());
  });
  it("should verify scalar multiple pairing", async () => {
    assert(await blsTest.pairingCheckMult.call());
  });
  it("should add points correctly", async () => {
    assert(await blsTest.addTest.call());
  })
  it("should do scalar multiplication correctly", async () => {
    assert(await blsTest.scalarTest.call());
  })
  it("should verify a simple signature correctly", async () => {
    assert(await blsTest.testSignature.call([12345,54321,10101,20202,30303]));
  })
  it("should sum points correctly", async () => {
    assert(await blsTest.testSumPoints.call());
  })
})
