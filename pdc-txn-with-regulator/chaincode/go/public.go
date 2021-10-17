package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (c *Contract) CreatePublicAsset(
	ctx contractapi.TransactionContextInterface,
	assetId string,
	name string,
) (*AssetPublic, error) {

	mspIdClient, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get msp assetId.\n%v", err)
	}

	owner64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get id.\n%v", err)
	}

	ownerBytes, _ := b64.StdEncoding.DecodeString(owner64)
	owner := string(ownerBytes)

	bytes, err := ctx.GetStub().GetState(assetId)
	if err != nil {
		return nil, fmt.Errorf("failed GetState for assetId: %s", assetId)
	}

	if len(bytes) != 0 {
		return nil, fmt.Errorf("asset already exist. assetId: %s", assetId)
	}

	assetPub := AssetPublic{
		ObjectType: assetType,
		AssetId:    assetId,
		Name:       name,
		Owner:      owner,
		OwnerOrg:   mspIdClient,
		Available:  true,
	}

	bytes, err = json.Marshal(assetPub)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal asset json. Asset: %v", assetPub)
	}

	if err := ctx.GetStub().PutState(assetId, bytes); err != nil {
		return nil, fmt.Errorf("failed PutState for Id: %s", assetId)
	}

	return &assetPub, nil
}

func (c *Contract) ReadPublicAsset(
	ctx contractapi.TransactionContextInterface,
	assetId string,
) (*AssetPublic, error) {

	bytes, err := ctx.GetStub().GetState(assetId)
	if err != nil {
		return nil, fmt.Errorf("failed GetState for assetId: %s", assetId)
	}

	if len(bytes) == 0 {
		return nil, fmt.Errorf("asset does not exist. assetId: %s", assetId)
	}

	var assetPub AssetPublic
	if err = json.Unmarshal(bytes, &assetPub); err != nil {
		return nil, errors.New("failed to unmarshal asset object")
	}

	return &assetPub, nil
}
