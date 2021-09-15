# telness-mananger
Manages the subscription(Create, Update and Get)

## Description
The Application exposes rest api's for creating, updataing subscription and  finding existing subscription in system.

The routes for application includes:
-----------------------------------------------

* Health: "/api/subscription/health"
* FindSubscription: "/api/subscription/{msidn}"
* CreateSubscription: "/api/subscription"
* UpdateSubscription: "/api/subscription"

msidn: define your subscription unique id.

Note: CreateSubscription & UpdateSubscription take json data to creat and update subscription

The URLS the application supports :
------------------------------------
* [Health](http://localhost:9000/api/subscription/health) 
* [FindSubscription](http://localhost:9000/api/subscription/{msidn})
* [CreateSubscription](http://localhost:9000/api/subscription)
* [UpdateSubscription](http://localhost:9000/api/subscription)

Note: Port is 8080 when using docker, else port is set to 9000 in .env file(when port cannot be accessed from env file, then default port is 8080).

## Application has:

- Go 1.16
- Makefile
- Dockerfile
- Docker-compose

## Running the application

* To build the application:

```bash
    make build
    cd bin/
    ./telness-manager
```

* To run test:
```bash
    make test
```


* To run integration test:
```bash
    make integration-test
```

* To run the application inside container:
```bash
    make up

    curl -X GET http://localhost:8080/api/subscription/health
    ------------------------------------------------------------------------
    curl -X POST http://localhost:8080/api/subscription -d '{"msidn": "c019ecde-17cb-4ef8-8a7d-85937a9250ed", "activate_at": "2021-09-13", "sub_type": "pbx", "status": "pending"}'
    ------------------------------------------------------------------------
    curl -X PATCH  http://localhost:8080/api/subscription -d '{"msidn": "c019ecde-17cb-4ef8-8a7d-85937a9250ed", "activate_at": "2021-09-15", "sub_type": "pbx", "status": "activated"}'
    ------------------------------------------------------------------------
    curl -X GET http://localhost:8080/api/subscription/c019ecde-17cb-4ef8-8a7d-85937a9250ed

```