### What is purpose of the task
- provide REST API for receiving data from MongoDB table
- implement service with ready in ready for production state, next things should be covered
    - service delivery
    - monitoring
    - authorization

### What is not purpose of the task
- data validation on loading into MongoDB
- covering complex use cases for data receiving: f.e. receive records where number > 5
- optimization queries in database

### Quick description:
- HTTP server for receiving data from MongoDB table.
- Light-weighted docker image ~15mb
- Provides /health and /metrics endpoints for monitoring

### How to run
- `make start` - starts web server and mongoDB inside docker. Data is being seed from mongo-seed/init.json file. For
development needs, mongoDB doesn't have persistent storage. Data will be reset by data from the seed file once
run `make start` command.

- `make stop` - stops the application and all containers

### How to use

With default params, service should respond to the requests
`curl "user:password@localhost:8080/api/v1/trivia?limit=1&found=true&number=60&type=trivia"`

Params:
- limit - count of responses
- found - true/false
- number - number
- type - string 

All params are optional.

### Delivery

- Service delivered in Docker container.
- Service can be configured via env variable, check .env file
- `/health` endpoint should be used to check if service is healthy
- Only `/api/v1/trivia` endpoint should be exposed to the unprotected network

### Monitoring

`/metrics` endpoint provides prometheus metrics for monitoring the application

It's recommended to add alerts for general metrics like CPU/memory along with:

- `http_requests_bucket` with status code > 2XX
- `http_requests_bucket` if response time is more than 100ms

### What should be improved
- use identity and access management solution instead of basicAuthorization
- improve test coverage
- add metrics for database communication: response time, number of requests etc.
- some DI should be used if number of services grows# redis_channels_playground
