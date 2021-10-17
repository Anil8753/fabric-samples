package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func (c *Contract) VerifyReceipt(
	ctx contractapi.TransactionContextInterface,
	orderId string,
) (bool, error) {

	mspId, err := ctx.GetClientIdentity().GetMSPID()
	if err != nil {
		return false, fmt.Errorf("failed to get msp id.\n%v", err)
	}

	if mspId != RegulatorMsp {
		return false, fmt.Errorf("only org2 is allowed to verify the order. passed org: %s", mspId)
	}

	collection := buildCollectionName(RegulatorMsp)

	// Order data hash
	order, err := c.GetOrder(ctx, orderId)
	if err != nil {
		return false, fmt.Errorf("failed to get private data.\n%v", err)
	}

	b, err := json.Marshal(order)
	if err != nil {
		return false, fmt.Errorf("failed to unmarshal order object.\n%v", err)
	}

	if len(b) == 0 {
		return false, fmt.Errorf("key not found.\n%v", err)
	}

	hash := sha256.New()
	hash.Write(b)
	hashCalculated := hash.Sum(nil)

	// ledger hash
	hashLedger, err := ctx.GetStub().GetPrivateDataHash(collection, orderId)
	if err != nil {
		return false, fmt.Errorf("failed to get the ledger hash for orderId: %s.\n%v", orderId, err)
	}

	if hashLedger == nil {
		return false, fmt.Errorf("asset private properties hash does not exist: %s", orderId)
	}

	// verify that the hash of the order object hash matches on-chain hash
	if !bytes.Equal(hashCalculated, hashLedger) {
		return false, fmt.Errorf("hash %x for passed orderId %s does not match on-chain hash %x. \n%s",
			hashCalculated,
			orderId,
			hashLedger,
			string(b),
		)
	}

	return true, nil
}
