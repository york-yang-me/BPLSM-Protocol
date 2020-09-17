let TestBPLSM = artifacts.require('BPLSM');

contract('BPLSM', (accounts) => {
    // the value of com, declares, seeds and H selected may need to be more precise data
    // you change these contents after you read my code
    it('should verify a valid SumCom', () => {
        return TestBPLSM.deployed().then((instance) => {
            let com = [12234,21231];
            let declares = [213,456,4546];
            let seeds = [12,121,1231];
            let H = 1231233323;
            return instance.verifyBPLSMSumCommit(com, declares, seeds, H);
        });
    });
    // the value of signs, R, pb_key and sumCom selected may need to be more precise data
    // you change these contents after you read my code
    it('should verify a valid BPLSM_SchnorrAggregateSign', () => {
        return TestBPLSM.deployed().then((instance) => {
            let signs = 2132135;
            let R = [1,2];
            let pb_key = [1,2];
            let sumCom = 112124524242424;
            return instance.verifyBPLSMSchnorr(signs, R, pb_key, sumCom);
        });
    });
});
