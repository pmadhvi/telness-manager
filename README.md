# telness-mananger
Manages the subscription(Create, Update and Get, and Update of status and activation_date )

## Description
The Application exposes rest api's for creating, updataing subscription and  finding existing subscription in system.

The routes for application includes:
-----------------------------------------------

* Health: "/api/subscription/health"
* FindSubscription: "/api/subscription/msisdn/{msisdn}"
* CreateSubscription: "/api/subscription"
* UpdateSubscription: "/api/subscription"
* UpdateStatusSubscription: "/api/subscription/update-subscription/msisdn/{msisdn}/status/{status}"
* UpdateActivateDate: "/api/subscription/update-activation-date/msisdn/{msisdn}/date/{date}"

msisdn: define your subscription unique number/phone number.
date: string value of future date

Note: CreateSubscription & UpdateSubscription take json data to create and update subscription

The URLS the application supports:
------------------------------------
* [Health](http://localhost:9000/api/subscription/health) 
* [FindSubscription](http://localhost:9000/api/subscription/msisdn/{msisdn})
* [CreateSubscription](http://localhost:9000/api/subscription)
* [UpdateSubscription](http://localhost:9000/api/subscription)
* [UpdateStatusSubscription](http://localhost:9000/api/subscription/update-subscription/msisdn/{msisdn}/status/{status})
* [UpdateActivateDate](http://localhost:9000/api/subscription/update-activation-date/msisdn/{msisdn}/date/{date})

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

    Run the curl request on another terminal(Note currently database is empty, so first run create request to create atleast one subscription):

    [Health check]: 
    curl -X GET http://localhost:8080/api/subscription/health

    ------------------------------------------------------------------------
    [Create Request]: 
    curl -X POST http://localhost:8080/api/subscription -d '{"msisdn": "+46107500500", "activate_at": "2021-10-13", "sub_type": "pbx", "status": "pending"}'

    ------------------------------------------------------------------------
    [Update Request]:
    curl -X PATCH  http://localhost:8080/api/subscription -d '{"msisdn": "+46107500500", "activate_at": "2021-10-15", "sub_type": "pbx", "status": "activated"}'

    ------------------------------------------------------------------------
    [Find Request]:
    curl -X GET http://localhost:8080/api/subscription/msisdn/+46107500500

    ------------------------------------------------------------------------
    [Update Status Request]: 
    status shoule be one of these 'cancelled', 'pending', 'activated', 'paused'
    curl -X PATCH http://localhost:8080/api/subscription/update-subscription/msisdn/+46107500500/status/cancelled

    ------------------------------------------------------------------------
    [Update activation date Request]:
    curl -X PATCH http://localhost:8080/api/subscription/update-activation-date/msisdn/+46107500500/date/2021-10-19

    ------------------------------------------------------------------------
```