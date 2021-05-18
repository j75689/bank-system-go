# bank-system-go

## Test
#### Create User
```bash
curl -X POST -H 'content-type: application/json' localhost:8080/api/v1/register -d '{"name": "test", "account": "test1", "password":"123456"}' -i
```
#### Login
```bash
curl -X POST -H 'content-type: application/json' localhost:8080/api/v1/login -d '{"account":"test1","password":"123456"}' -i
```

#### Create Wallet
```bash
JWT_TOEKN=...
curl -X POST -H 'content-type: application/json' -H "authentication: Bearer ${JWT_TOKEN}" localhost:8080/api/v1/wallet -d '{"type":1,"currency_id":1}' -i
```

#### List Wallet
```bash
JWT_TOEKN=...
curl -X GET -H 'content-type: application/json' -H "authentication: Bearer ${JWT_TOKEN}" localhost:8080/api/v1/wallets -d '{"pagination":{"page":0,"per_page":10},"sort":[{"sort_field":"id","sort_order":"DESC"}]}' -i
```

#### Update Wallet Balance
```bash
JWT_TOEKN=...
ACCOUNT=...
# deposit
curl -X POST -H 'content-type: application/json' -H "authentication: Bearer ${JWT_TOKEN}" localhost:8080/api/v1/wallet/balance -d "{\"type\":1,\"account_number\":\"${ACCOUNT}\",\"amount\":100}" -i
# withdraw
curl -X POST -H 'content-type: application/json' -H "authentication: Bearer ${JWT_TOKEN}" localhost:8080/api/v1/wallet/balance -d "{\"type\":2,\"account_number\":\"${ACCOUNT}\",\"amount\":100}" -i
```

#### List Transaction
```bash
JWT_TOEKN=...
curl -X GET -H 'content-type: application/json' -H "authentication: Bearer ${JWT_TOKEN}" localhost:8080/api/v1/transactions -d "{\"account_number\":\"${ACCOUNT}\",\"pagination\":{\"page\":0,\"per_page\":50},\"created_at_lte\":1621351737}" -i
```