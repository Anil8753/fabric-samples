package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-chaincode-go/pkg/statebased"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

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

func getTransientInput(ctx contractapi.TransactionContextInterface) (*TransientInput, error) {

	transientMap, err := ctx.GetStub().GetTransient()
	if err != nil {
		return nil, fmt.Errorf("failed to get transoent data.\n%v", err)
	}

	data, ok := transientMap["data"]
	if !ok {
		return nil, fmt.Errorf("failed to get 'data' field from the transient data.\n%v", transientMap)
	}

	var input TransientInput
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("failed to parse transoent data.\n%v", err)
	}

	return &input, nil
}

func markAssetSold(
	ctx contractapi.TransactionContextInterface,
	assetPub *AssetPublic,
) error {

	assetPub.Available = false

	bytes, err := json.Marshal(assetPub)
	if err != nil {
		return fmt.Errorf("failed to marshal asset json. Asset: %v", assetPub)
	}

	if err := ctx.GetStub().PutState(assetPub.AssetId, bytes); err != nil {
		return fmt.Errorf("failed PutState for Id: %s", assetPub.AssetId)
	}

	return nil
}
