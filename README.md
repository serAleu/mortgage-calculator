
Примеры запросов для расчета ипотеки:

    1: curl -X POST http://localhost:8282/execute   -H "Content-Type: application/json"   -d '{
        "object_cost": 500870000,
        "initial_payment": 1000086500,
        "months": 400,
        "program": {
        "military": true
        }
    }'
    2: curl -X POST http://localhost:8282/execute \
    -H "Content-Type: application/json" \
    -d '{
        "object_cost": 5000000,
        "initial_payment": 1000000,
        "months": 240,
        "program": {
        "salary": true
        }
    }'

Получить значения из кэша:

curl http://localhost:8282/cache