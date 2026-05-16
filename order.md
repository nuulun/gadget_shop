curl -X POST http://34.122.34.46/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"password"}'


curl -X POST http://34.122.34.46/api/orders \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer " \
  -d '{"user_id": 1, "items": [{"product_id": 1, "quantity": 1}]}'


  ab -t 60 -c 50 \
  -p /dev/stdin \
  -T "application/json" \
  -H "Authorization: Bearer " \
  http://34.122.34.46/api/orders \
  <<< '{"user_id": 1, "items": [{"product_id": 1, "quantity": 1}]}'

