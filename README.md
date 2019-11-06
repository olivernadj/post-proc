# post-proc
This is a sample server processing and saving incoming requests


## run
Make sure you have `docker` and `docker-compose` installed.
`$ make createdb`

- Swagger gui http://localhost:8080
- Grafana http://localhost:3000 user:admin, pass:5ecret
- MYSQL admin http://localhost:8081 user:root, pass:example

## test
`$ make test`

## Todo
The main goal of this test task is to develop the application for processing the incoming requests from the 3d-party providers.
The application must have an HTTP URL to receive incoming POST requests.
To receive the incoming POST requests the application must have an HTTP URL endpoint. 

### Requirements:
#### 1 Processing and saving incoming requests.

Example of the POST request: 
```
POST /your_url HTTP/1.1
Source-Type: client
Content-Length: 34
Host: 127.0.0.1
Content-Type: application/json

{"action": "save", "state": "new"}
```

Header `Source-Type` could be in 3 types (`client, server, payment`). This type probably could be expanded in the future.

The decision regarding database architecture and table structure is made to you. 

#### 2 Post-processing

- Every N minutes 10 new records must be marked as `processed`
- Every N minutes records marked as processed must change their state to `deleted` 
