package main

import (
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const (
	assetType    = "asset"
	RegulatorMsp = "org2-test-com"
)

type AssetPublic struct {
	ObjectType string `json:"objectType"`
	AssetId    string `json:"assetId"`
	Name       string `json:"name"`

	Owner    string `json:"owner"`
	OwnerOrg string `json:"ownerOrg"`

	Available bool `json:"available"`
}

type Order struct {
	OrderId string `json:"orderId"`
	AssetId string `json:"assetId"`

	Buyer     string `json:"buyer"`
	BuyerOrg  string `json:"buyerOrg"`
	Seller    string `json:"seller"`
	SellerOrg string `json:"sellerOrg"`

	Price string `json:"price"`
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
