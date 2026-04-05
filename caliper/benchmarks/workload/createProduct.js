'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class CreateProductWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this.workerIndex = workerIndex;
    }

    async submitTransaction() {
        this.txIndex++;
        const productId = `P_${this.workerIndex}_${this.txIndex}_${Date.now()}`;
        const location = `Farm-${this.workerIndex}`;
        const quality = 'Grade-A';

        const request = {
            contractId: 'agrifood',
            contractFunction: 'CreateProduct',
            invokerIdentity: 'admin-org1',
            contractArguments: [productId, location, quality],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new CreateProductWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
