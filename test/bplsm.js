let BPLSM = artifacts.require('BPLSM');

contract('BPLSM', (accounts) => {

    it('should verify a valid SumCom', () => {
        return BPLSM.deployed().then((instance) => {
            let com = [12234,21231];
            let declares = [213,456,4546];
            let seeds = [12,121,1231];
            let H = 1231233323;
            return instance.verifyBPLSMSumCommit(com, declares, seeds, H);
        });
    });
    it('should verify a valid BPLSM_SchnorrAggregateSign', () => {
        return BPLSM.deployed().then((instance) => {
            let signs = 2132135;
            let R = [1,2];
            let pb_key = [1,2];
            let sumCom = 112124524242424;
            return instance.verifyBPLSMSchnorr(signs, R, pb_key, sumCom);
        });
    });
});
