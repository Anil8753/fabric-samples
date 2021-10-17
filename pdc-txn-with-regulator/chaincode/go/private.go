package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (c *Contract) BuyAsset(
	ctx contractapi.TransactionContextInterface,
) (*Order, error) {

	buyerOrg, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get msp id.\n%v", err)
	}

	buyer64, err := ctx.GetClientIdentity().GetID()
	if err != nil {
		return nil, fmt.Errorf("failed to get id.\n%v", err)
	}

	buyerBytes, _ := b64.StdEncoding.DecodeString(buyer64)
	buyer := string(buyerBytes)

	// get transient input
	transientInput, err := getTransientInput(ctx)
	if err != nil {
		return nil, err
	}

	assetPub, err := c.ReadPublicAsset(ctx, transientInput.AssetId)
	if err != nil {
		return nil, fmt.Errorf("failed to get Asset\n%v", err)
	}

	if !assetPub.Available {
		return nil, fmt.Errorf("asset already soldout. assetId: %s\n%v", transientInput.AssetId, err)
	}

	if assetPub.Owner == buyer {
		return nil, fmt.Errorf("an owner cannot buy himself. identity: %s", buyer)
	}

	collectionRegulator := buildCollectionName(RegulatorMsp)

	// endorsement policy
	// if err := setPrivateStateBasedEndorsement(ctx, collection, transientInput.OrderId, mspIdClient, assetPub.Owner); err != nil {
	// 	return nil, err
	// }

	order := Order{
		OrderId: transientInput.OrderId,
		AssetId: transientInput.AssetId,

		Buyer:     buyer,
		BuyerOrg:  buyerOrg,
		Seller:    assetPub.Owner,
		SellerOrg: assetPub.OwnerOrg,

		Price: transientInput.Price,
	}

	bytes, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data.\n%v", err)
	}

	if err := ctx.GetStub().PutPrivateData(collectionRegulator, order.OrderId, bytes); err != nil {
		return nil, fmt.Errorf("failed to put private data in regulator collection.\n%v", err)
	}

	// Make the Asset soldout
	if err := markAssetSold(ctx, assetPub); err != nil {
		// Todo: make the order cancelled
		return nil, fmt.Errorf("failed to mark asset sold.\n%v", err)
	}

	// Change the current Owner of the asset
	assetPub.Owner = buyer
	assetPub.OwnerOrg = buyerOrg
	bytes, err = json.Marshal(assetPub)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data.\n%v", err)
	}

	if err := ctx.GetStub().PutState(assetPub.AssetId, bytes); err != nil {
		return nil, fmt.Errorf("failed to put the update asset with new owner.\n%v", err)
	}
	return &order, nil
}

func (c *Contract) GetOrder(
	ctx contractapi.TransactionContextInterface,
	orderId string,
) (*Order, error) {

	mspId, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return nil, fmt.Errorf("failed to get msp id.\n%v", err)
	}

	collection := buildCollectionName(RegulatorMsp)

	b, err := ctx.GetStub().GetPrivateData(collection, orderId)
	if err != nil {
		return nil, fmt.Errorf("failed to get private data.\n%v", err)
	}

	if len(b) == 0 {
		return nil, fmt.Errorf("key not found.\n Collection: %s, ClientMspId: %s", collection, mspId)
	}

	var order Order
	if err := json.Unmarshal(b, &order); err != nil {
		return nil, fmt.Errorf("failed to unmarshal the bytes. %v", err)
	}

	return &order, nil
}
