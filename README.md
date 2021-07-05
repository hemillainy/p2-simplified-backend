# p2-simplified-backend

# Build
```bash
docker-compose build
```

# Run
```bash
docker-compose up
```

# Make a transfer
```
curl --location --request POST 'http://localhost:8080/transaction' \
--header 'Content-Type: application/json' \
--data-raw '{
    "value": 100,
    "payer": "4bfa93a9-417a-40d5-b18b-2335cff9aeac",
    "payee": "bee22ab3-9ba7-4d07-b372-70d5d12c56ec"
}'
```