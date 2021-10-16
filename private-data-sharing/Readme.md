# Sharing implicit data collection data with other organization
Implicit data collecton follows the naming convention `_implicit_org_%s` where `%s` is the MSP Id of the Org.
We can use the ChaincodeStubInterface API mentioned below to write the private details directly to the implicit private 
data collections of the desired orgs

```
	// PutPrivateData puts the specified `key` and `value` into the transaction's
	// private writeset. Note that only hash of the private writeset goes into the
	// transaction proposal response (which is sent to the client who issued the
	// transaction) and the actual private writeset gets temporarily stored in a
	// transient store. PutPrivateData doesn't effect the `collection` until the
	// transaction is validated and successfully committed. Simple keys must not
	// be an empty string and must not start with a null character (0x00) in order
	// to avoid range query collisions with composite keys, which internally get
	// prefixed with 0x00 as composite key namespace. In addition, if using
	// CouchDB, keys can only contain valid UTF-8 strings and cannot begin with an
	// an underscore ("_").
	PutPrivateData(collection string, key string, value []byte) error
```

Basically it is a simplified version of 
asset-transfer-secured-agreement (https://github.com/hyperledger/fabric-samples/tree/master/asset-transfer-secured-agreement/chaincode-go) chaincode in fabric sample.
Because the purpose is only show how can we share the implicit private data collection
with other orgs, many checks like only current owner can transfer the asset to other
org are removed. 


#
# Make network up and running

### Copy the chaincode
Copy the chaincode into the vars/chaincode directory

```
mkdir -p ./vars/chaincode/pdc && cp -rf ./chaincode/go ./vars/chaincode/pdc/
```

### Create network
```
./minifab up -s couchdb -n pdc -p '"InitLedger"' -e true -r true -o org0.test.com
```

### Create Public Asset
```
./minifab invoke -p '"CreatePublicAsset", "asset_001", "First Asset"' -o org0.test.com
```

### Read Public Asset
By Org0
```
./minifab query -p '"ReadPublicAsset", "asset_001"' -o org0.test.com
```

By Org1
```
./minifab query -p '"ReadPublicAsset", "asset_001"' -o org1.test.com
```

By Org2
```
./minifab query -p '"ReadPublicAsset", "asset_001"' -o org2.test.com
```

#
# Try buy the asset

#### If org0-test-com tries to buy it should fail as owner cannot buy his asset himself
Pass the input data as transient field as compared to the arguments list
```
export KEYVALUE=$(echo -n '{ "assetId":"asset_001", "orderId":"order_001", "price":"1000" }' | base64 | tr -d \\n)
./minifab invoke -p '"BuyAsset"' -t '{ "data": "'$KEYVALUE'"}' -o org0.test.com
```

#### If org1-test-com tries to buy it should sucessful
Pass the input data as transient field as compared to the arguments list
```
export KEYVALUE=$(echo -n '{ "assetId":"asset_001", "orderId":"order_001", "price":"1000" }' | base64 | tr -d \\n)
./minifab invoke -p '"BuyAsset"' -t '{ "data": "'$KEYVALUE'"}' -o org1.test.com
```

#### If org2-test-com tries to buy it should fail as asset is already sold
Pass the input data as transient field as compared to the arguments list
```
export KEYVALUE=$(echo -n '{ "assetId":"asset_001", "orderId":"order_001", "price":"2000" }' | base64 | tr -d \\n)
./minifab invoke -p '"BuyAsset"' -t '{ "data": "'$KEYVALUE'"}' -o org2.test.com
```

### Query private transactions
```
./minifab query -p '"GetOrder", "order_001"' -o org1.test.com
```

### Check current owner of the asset, it should be org1-test-com
```
./minifab query -p '"ReadPublicAsset", "asset_001"' -o org0.test.com
```

#
# Bring network down and cleanup

### Bring down network
```
./minifab down
```

### Clean everything
```
./minifab cleanup
```