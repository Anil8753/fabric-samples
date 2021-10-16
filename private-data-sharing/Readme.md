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
