package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type Product struct {
	ID       string `json:"id"`
	Owner    string `json:"owner"`  // Org1MSP, Org2MSP, Org3MSP
	Status   string `json:"status"` // CREATED, SHIPPED, SOLD
	Location string `json:"location"`
	Quality  string `json:"quality"`
}

// CreateProduct - Only Org1 (Farmer) can create
func (s *SmartContract) CreateProduct(ctx contractapi.TransactionContextInterface, id, location, quality string) error {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil || mspID != "Org1MSP" {
		return fmt.Errorf("only Org1 (Farmer) can create products")
	}

	exists, _ := s.ProductExists(ctx, id)
	if exists {
		return fmt.Errorf("product %s already exists", id)
	}

	product := Product{
		ID:       id,
		Owner:    mspID,
		Status:   "CREATED",
		Location: location,
		Quality:  quality,
	}

	productJSON, _ := json.Marshal(product)
	return ctx.GetStub().PutState(id, productJSON)
}

// UpdateProductStatus - Only Org2 (Retailer) can update, and only if current owner is Org1MSP
func (s *SmartContract) UpdateProductStatus(ctx contractapi.TransactionContextInterface, id, status, location string) error {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil || mspID != "Org2MSP" {
		return fmt.Errorf("only Org2 (Retailer) can update product status")
	}

	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil || productJSON == nil {
		return fmt.Errorf("product %s does not exist", id)
	}

	var product Product
	_ = json.Unmarshal(productJSON, &product)

	if product.Owner != "Org1MSP" {
		return fmt.Errorf("product %s is not owned by Org1; cannot be updated by Org2", id)
	}

	product.Status = status
	product.Location = location
	product.Owner = mspID // Transfer ownership to Org2

	updatedJSON, _ := json.Marshal(product)
	return ctx.GetStub().PutState(id, updatedJSON)
}

// PurchaseProduct - Only Org3 (Customer) can purchase, and only if owned by Org2MSP
func (s *SmartContract) PurchaseProduct(ctx contractapi.TransactionContextInterface, id string) error {
	mspID, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil || mspID != "Org3MSP" {
		return fmt.Errorf("only Org3 (Customer) can purchase products")
	}

	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil || productJSON == nil {
		return fmt.Errorf("product %s does not exist", id)
	}

	var product Product
	_ = json.Unmarshal(productJSON, &product)

	if product.Owner != "Org2MSP" {
		return fmt.Errorf("product %s is not owned by Org2; cannot be sold to Org3", id)
	}

	product.Status = "SOLD"
	product.Owner = mspID

	updatedJSON, _ := json.Marshal(product)
	return ctx.GetStub().PutState(id, updatedJSON)
}

// ProductExists utility function
func (s *SmartContract) ProductExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	data, err := ctx.GetStub().GetState(id)
	return data != nil, err
}

// GetProduct - Anyone can read product by ID
func (s *SmartContract) GetProduct(ctx contractapi.TransactionContextInterface, id string) (*Product, error) {
	productJSON, err := ctx.GetStub().GetState(id)
	if err != nil || productJSON == nil {
		return nil, fmt.Errorf("product %s not found", id)
	}

	var product Product
	_ = json.Unmarshal(productJSON, &product)
	return &product, nil
}

// GetAllProducts - Return all products (open query for simplicity)
func (s *SmartContract) GetAllProducts(ctx contractapi.TransactionContextInterface) ([]*Product, error) {
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer iterator.Close()

	var products []*Product
	for iterator.HasNext() {
		item, err := iterator.Next()
		if err != nil {
			return nil, err
		}
		var product Product
		_ = json.Unmarshal(item.Value, &product)
		products = append(products, &product)
	}

	return products, nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(SmartContract))
	if err != nil {
		panic("Error creating chaincode: " + err.Error())
	}
	if err := chaincode.Start(); err != nil {
		panic("Error starting chaincode: " + err.Error())
	}
}
