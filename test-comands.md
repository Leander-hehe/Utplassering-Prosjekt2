# Test commands
- curl http://localhost:8080/
- curl http://localhost:8080/1
- curl http://localhost:8080/healthcheck
- curl http://localhost:8080/ -i -H "Content-Type: application/json" -X "POST" --data '{"id": "2", "name": "B"}'
- curl http://localhost:8080/1 -i -H "Content-Type: application/json" -X "DELETE"
- curl http://localhost:8080/1 -i -H "Content-Type: application/json" -X "PUT" --data '{"id": "1", "name": "Z"}'
- curl http://localhost:8080/1 -i -H "Content-Type: application/json" -X "PATCH" --data '[["Name", "Z"]]'


## With X-API-Key
- curl http://localhost:8080/ -H "X-API-Key: test123"
- curl http://localhost:8080/1 -H "X-API-Key: test123"
- curl http://localhost:8080/healthcheck -H "X-API-Key: test123"
- curl http://localhost:8080/ -i -H "Content-Type: application/json" -X "POST" --data '{"id": "2", "name": "B"}' -H "X-API-Key: test123"
- curl http://localhost:8080/1 -i -H "Content-Type: application/json" -X "DELETE" -H "X-API-Key: test123"
- curl http://localhost:8080/1 -i -H "Content-Type: application/json" -X "PUT" --data '{"id": "1", "name": "Z"}' -H "X-API-Key: test123"
- curl http://localhost:8080/1 -i -H "Content-Type: application/json" -X "PATCH" --data '[["Name", "Z"]]' -H "X-API-Key: test123"
