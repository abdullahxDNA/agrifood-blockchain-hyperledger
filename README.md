# Blockchain Implementation for Agricultural Food Supply Chain using Hyperledger Fabric

> Undergraduate Thesis — International Islamic University Chittagong (IIUC)  
> Department of Computer Science and Engineering

## Overview

The agricultural food supply chain in Bangladesh suffers from a lack of transparency, inefficiency, and consumer distrust — leading to food fraud, contamination, and economic losses for farmers. This project presents a **decentralized traceability system** built on **Hyperledger Fabric**, a permissioned blockchain framework, to track agricultural products from farm to table.

The system provides a secure, immutable, and transparent ledger that enforces strict role-based access control across all stakeholders.

---

## Network Architecture

A 3-organization Hyperledger Fabric network representing the key stakeholders:

| Organization | Role | Chaincode Permission |
|---|---|---|
| **Org1** | Farmer | `CreateProduct` |
| **Org2** | Retailer | `UpdateProductStatus` |
| **Org3** | Customer | `PurchaseProduct` |

All organizations share a single channel (`mychannel`) with a single orderer using the **Raft consensus algorithm**.

![Network Architecture](docs/figures/Figure%206%20%20Network%20Architecture.jpg)

---

## Supply Chain Workflow

```
Farmer (Org1)  →  Retailer (Org2)  →  Customer (Org3)
CreateProduct     UpdateProductStatus   PurchaseProduct
Status: CREATED   Status: SHIPPED       Status: SOLD
Owner: Org1MSP    Owner: Org2MSP        Owner: Org3MSP
```

![Workflow](docs/figures/Figure%207%20%20Agri-Food%20Supply%20Chain%20Workflow.jpg)

---

## Smart Contract (Chaincode)

Written in **Go** using the `fabric-contract-api-go` library. Located at `chaincode/agrifood/agrifood.go`.

### Product Structure
```go
type Product struct {
    ID       string `json:"id"`
    Owner    string `json:"owner"`    // Org1MSP, Org2MSP, Org3MSP
    Status   string `json:"status"`   // CREATED, SHIPPED, SOLD
    Location string `json:"location"`
    Quality  string `json:"quality"`
}
```

### Functions

| Function | Caller | Description |
|---|---|---|
| `CreateProduct` | Org1 (Farmer) | Introduces a new product to the ledger |
| `UpdateProductStatus` | Org2 (Retailer) | Updates status and transfers ownership from Org1 |
| `PurchaseProduct` | Org3 (Customer) | Marks product as SOLD, transfers ownership from Org2 |
| `GetProduct` | Any | Query a product by ID |
| `GetAllProducts` | Any | Query all products on the ledger |

---

## Tech Stack

- **Blockchain**: Hyperledger Fabric v2.5
- **Chaincode Language**: Go
- **Consensus**: Raft
- **State Database**: LevelDB
- **Containerization**: Docker & Docker Compose
- **OS**: Ubuntu (VMware Virtual Machine)
- **Explorer**: Hyperledger Explorer

---

## Prerequisites

- Docker & Docker Compose
- Git
- curl

## Setup & Run

**1. Install Fabric binaries and Docker images:**
```bash
./install-fabric.sh docker binary
```

**2. Start the test network with 3 organizations:**
```bash
cd test-network
./network.sh up createChannel -c mychannel -ca
./network.sh addOrg3
```

**3. Deploy the agrifood chaincode:**
```bash
./network.sh deployCC -ccn agrifood -ccp ../chaincode/agrifood -ccl go
```

**4. Set environment for Org1 (Farmer) and create a product:**
```bash
export CORE_PEER_LOCALMSPID="Org1MSP"
# (set peer TLS and address env vars for Org1)

peer chaincode invoke \
  -C mychannel -n agrifood \
  -c '{"function":"CreateProduct","Args":["P001","Farm A","Grade A"]}'
```

**5. Org2 (Retailer) updates status:**
```bash
# (set env vars for Org2)
peer chaincode invoke \
  -C mychannel -n agrifood \
  -c '{"function":"UpdateProductStatus","Args":["P001","SHIPPED","Warehouse"]}'
```

**6. Org3 (Customer) purchases:**
```bash
# (set env vars for Org3)
peer chaincode invoke \
  -C mychannel -n agrifood \
  -c '{"function":"PurchaseProduct","Args":["P001"]}'
```

**7. Query final state:**
```bash
peer chaincode query \
  -C mychannel -n agrifood \
  -c '{"function":"GetProduct","Args":["P001"]}'
```

**8. Tear down the network:**
```bash
./network.sh down
```

---

## Performance Results

Benchmarked using **Hyperledger Caliper** (1000 transactions per round):

| Operation | TPS | Avg Latency | Max Latency | Success Rate |
|---|---|---|---|---|
| `CreateProduct` (Write) | **93.2** | 270 ms | 860 ms | 1000/1000 (100%) |
| `GetAllProducts` (Read) | **55.2** | 6920 ms | 9360 ms | 1000/1000 (100%) |

Compared against other platforms:

| Platform | TPS | Latency |
|---|---|---|
| **This Project (Fabric + Raft)** | **93.2** | **270 ms** |
| Hyperledger Sawtooth (PoET) | 45 | 203 ms |
| Ethereum (PoW) | 30 | 214 ms |

---

## Project Documents

| Document | Description |
|---|---|
| [IEEE Paper](docs/IEEE%20format(Blockchain_Implementation_for_Agricultural_Food_Supply_Chain_using_Hyperledger_Fabric).pdf) | Full research paper in IEEE format |
| [Thesis Report](docs/thesis-report-final.pdf) | Final thesis report (B.Sc.) |
| [Presentation](docs/Defence_ppt.pptx) | Defense presentation slides |
| [Figures](docs/figures/) | All architecture and result diagrams |

---

## Authors

| Name | Student ID |
|---|---|
| **Ashraf Ali Rakib** | C213090 |
| **Miraj Mahmud** | C213062 |
| **Md Abdullah** | C213042 |

**Supervisor:** Abdullahil Kafi — Assistant Professor, Department of CSE, IIUC

B.Sc. in Computer Science and Engineering  
International Islamic University Chittagong (IIUC)  
Spring 2025
