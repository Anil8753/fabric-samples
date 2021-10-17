# Sharing implicit data collection data with other organization
Implicit data collecton follows the naming convention `_implicit_org_%s` where `%s` is the MSP Id of the Org.
We can use the ChaincodeStubInterface API mentioned below to write the private details directly to the implicit private 
data collections of the desired orgs

```
	PutPrivateData(collection string, key string, value []byte) error
```

In this sample, Asset transaction receipt (order) is written into the implicit data collection of 
the regulator (org2) and org2 can verify the saved order id (receipt) in the private data collection


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

### Query private data transactions
As order data is saved in implicit data collection of `org2`, only `org2` can query
```
./minifab query -p '"GetOrder", "order_001"' -o org2.test.com
```

### Query private data verification
As transaction (order) data is saved in implicit data collection of `org2` only `org2` can verify in this sample.
In practice you can share the data off-chain with orgs and other orgs also can verify the same.
But `VerifyReceipt` function functionality works only for `org2`.
```
./minifab query -p '"VerifyReceipt", "order_001"' -o org2.test.com
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