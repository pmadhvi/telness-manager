# iban-validator
Validates IBAN for different countries

IBAN Validator API
==================================

## Description
The Application checks if provided iban is a valid iban or not

The routes for application includes:
-----------------------------------------------

* Health: "/api/iban/validate/health"
* ValidateIban: "/api/iban/validate/{iban}"

iban: define your country iban

The URLS the application supports :
------------------------------------
* [Health](http://localhost:9000/api/iban/validate/health) 
* [ValidateIban](http://localhost:9000//api/iban/validate/{iban})

Note: Port is 8080 when using docker, else port is set to 9000 in .env file(when port cannot be accessed from env file, then default port is 8080).

## Application has:

- Go 1.16
- Makefile
- Dockerfile
- Docker-compose

## Running the application

* To build the application on mac-osx:

```bash
    make build
    cd bin/
    ./iban-validator
```

* To build the application on linux-32:

```bash
    make build-linux-32
    cd bin/
    ./iban-validator-linux-32
```

* To build the application on linux-64:

```bash
    make build-linux-64
    cd bin/
    ./iban-validator-linux-64
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
    docker-compose up
    curl http://localhost:8080/api/iban/validate/health
    curl http://localhost:8080/api/iban/validate/BA391290079401028494
```

## Iban Validation Reference: 
https://www.morfoedro.it/doc.php?n=219&lang=en#:~:text=The%20IBAN%20must%20have%20a,digits%20from%200%20to%209.