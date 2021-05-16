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