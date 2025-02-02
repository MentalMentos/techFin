TO RUN: make run

requests(raw body):
POST
http://localhost:8080/api/balance-update
{
    "user_id":1,
    "amount":123
}

POST
http://localhost:8080/api/transfer
{
    "from_user_id":1,
    "to_user_id":2,
    "amount":123
}

GET
http://localhost:8080/api/transactions
{
    "user_id":1
}
