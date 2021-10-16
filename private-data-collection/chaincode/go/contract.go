package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/pkg/statebased"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

const assetType = "asset"

type AssetPrivate struct {
	ObjectType string `json:"objectType"`
	Id         string `json:"id"`
	Name       string `json:"name"`
	Desc       string `json:"desc"`
}

type Contract struct {
	contractapi.Contract
}

func (c *Contract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	return nil
}

func (c *Contract) PutPrivateData(
	ctx contractapi.TransactionContextInterface,
	id string,
	name string,
	desc string,
) error {

	mspId, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return fmt.Errorf("failed to get msp id.\n%v", err)
	}

	collection := buildCollectionName(mspId)

	// endorsement policy
	if err := setPrivateStateBasedEndorsement(ctx, collection, id, mspId); err != nil {
		return fmt.Errorf("failed to set private data endorsement policy./n%v", err)
	}

	ap := AssetPrivate{
		ObjectType: assetType,
		Id:         id,
		Name:       name,
		Desc:       desc,
	}

	b, err := json.Marshal(ap)
	if err != nil {
		return fmt.Errorf("failed to marshal data.\n%v", err)
	}

	if err := ctx.GetStub().PutPrivateData(collection, id, b); err != nil {
		return fmt.Errorf("failed to put private data.\n%v", err)
	}

	return nil
}

func (c *Contract) GetPrivateData(
	ctx contractapi.TransactionContextInterface,
	id string,

) (*AssetPrivate, error) {

	mspId, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get msp id.\n%v", err)
	}

	collection := buildCollectionName(mspId)

	b, err := ctx.GetStub().GetPrivateData(collection, id)
	if err != nil {
		return nil, fmt.Errorf("failed to put private data.\n%v", err)
	}

	if len(b) == 0 {
		return nil, fmt.Errorf("key not found.\n%v", err)
	}

	var ap AssetPrivate
	if err := json.Unmarshal(b, &ap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the bytes. %v", err)
	}

	return &ap, nil
}

func buildCollectionName(clientOrgID string) string {
	return fmt.Sprintf("_implicit_org_%s", clientOrgID)
}

func setPrivateStateBasedEndorsement(
	ctx contractapi.TransactionContextInterface,
	collection string,
	id string,
	orgsToEndorse ...string,
) error {

	endorsementPolicy, err := statebased.NewStateEP(nil)
	if err != nil {
		return err
	}

	err = endorsementPolicy.AddOrgs(statebased.RoleTypeMember, orgsToEndorse...)
	if err != nil {
		return fmt.Errorf("failed to add org to endorsement policy: %v", err)
	}

	policy, err := endorsementPolicy.Policy()
	if err != nil {
		return fmt.Errorf("failed to create endorsement policy bytes from org: %v", err)
	}

	err = ctx.GetStub().SetPrivateDataValidationParameter(collection, id, policy)
	if err != nil {
		return fmt.Errorf("failed to set validation parameter on asset: %v", err)
	}

	return nil
}
