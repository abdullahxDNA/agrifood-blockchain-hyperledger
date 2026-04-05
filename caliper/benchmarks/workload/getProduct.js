'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class GetProductWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async submitTransaction() {
        const request = {
            contractId: 'agrifood',
            contractFunction: 'GetAllProducts',
            invokerIdentity: 'admin-org1',
            contractArguments: [],
            readOnly: true
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new GetProductWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
