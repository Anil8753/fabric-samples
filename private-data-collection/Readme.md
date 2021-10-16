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
```
./minifab invoke -p '"PutPrivateData", "001", "First", "This is first"' -o org0.test.com
```

```
./minifab invoke -p '"PutPrivateData", "002", "Second", "This is second"' -o org1.test.com
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