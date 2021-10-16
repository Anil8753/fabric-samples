package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const assetType = "asset"

type AssetPublic struct {
	ObjectType string `json:"objectType"`
	AssetId    string `json:"assetId"`
	Name       string `json:"name"`
	Owner      string `json:"owner"`
	Available  bool   `json:"available"`
}

type Order struct {
	OrderId string `json:"orderId"`
	AssetId string `json:"assetId"`
	Price   string `json:"price"`
}

type TransientInput struct {
	AssetId string `json:"assetId"`
	OrderId string `json:"orderId"`
	Price   string `json:"price"`
}

type Contract struct {
	contractapi.Contract
}

func (c *Contract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}
