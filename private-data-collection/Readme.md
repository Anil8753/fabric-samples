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

### Create Private Data
Pass the input data as transient field as compared to the arguments list
```
export KEYVALUE=$(echo -n '{ "id":"001", "name":"First", "desc":"This is first" }' | base64 | tr -d \\n)
./minifab invoke -p '"PutPrivateData"' -t '{ "data": "'$KEYVALUE'"}' -o org0.test.com
```

```
export KEYVALUE=$(echo -n '{ "id":"002", "name":"Second", "desc":"This is second" }' | base64 | tr -d \\n)
./minifab invoke -p '"PutPrivateData"' -t '{ "data": "'$KEYVALUE'"}' -o org1.test.com
```

### Read Private Data
```
./minifab query -p '"GetPrivateData", "001"' -o org0.test.com
```
```
./minifab query -p '"GetPrivateData", "002"' -o org1.test.com
```

#### Should fail
```
./minifab query -p '"GetPrivateData", "001"' -o org1.test.com
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