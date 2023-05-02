# Fizzbuzz

 [![Codacy Badge](https://app.codacy.com/project/badge/Grade/d097d0142e6043a3936879cd0433a696?refresh=1)](https://app.codacy.com/gh/arckadious/fizzbuzz/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)  [![Go Report Card](https://goreportcard.com/badge/github.com/arckadious/fizzbuzz?refresh=1)](https://goreportcard.com/report/github.com/arckadious/fizzbuzz)  [![codecov](https://codecov.io/gh/arckadious/fizzbuzz/branch/master/graph/badge.svg)](https://codecov.io/gh/arckadious/fizzbuzz)  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/arckadious/fizbuzz/LICENSE)  [![Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)

The original fizzbuzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by "fizz", all multiples of 5 by "buzz", and all multiples of 15 by "fizzbuzz".

Fizzbuzz is a web server. It has two REST API endpoints, one will return the list of numbers, and another one will show what the most frequent request has been.

______________

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

- [:arrows_clockwise: Requirements](#arrows_clockwise-requirements)
  - [Local Requirements](#local-requirements)
- [:package: Dependencies](#package-dependencies)
  - [Indirect Dependencies](#indirect-dependencies)
- [:vertical_traffic_light: Usages](#vertical_traffic_light-usages)
  - [Remote Uses](#remote-uses)
  - [Start containers](#start-containers)
  - [Local Installation](#local-installation)
  - [Run API Server (with air)](#run-api-server-with-air)
  - [Run API Server (without air)](#run-api-server-without-air)
  - [Run Tests](#run-tests)
- [:whale: Environment](#whale-environment)
- [:mag: Configuration](#mag-configuration)
- [:page_with_curl: Logs](#page_with_curl-logs)
  - [Applicative logs](#applicative-logs)
  - [Gin Framework Logs](#gin-framework-logs)
- [:link: Tools](#link-tools)
- [:trident: API Endpoints](#trident-api-endpoints)
  - [Main endpoint](#main-endpoint)
  - [Most request used](#most-request-used)
  - [Doc OpenAPI(swagger)](#doc-openapiswagger)
  - [:green_heart: Health check](#green_heart-health-check)
    - [Ping](#ping)
- [:file_folder: Workspace](#file_folder-workspace)
- [:briefcase: Author](#briefcase-author)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->
<!-- MSYS_NO_PATHCONV=1 docker run --rm -v "$(pwd)":/app peterdavehello/npm-doctoc doctoc /app/README.md -->

## :arrows_clockwise: Requirements

Golang 1.19+

### Local Requirements

- [docker v23.0.5](https://docs.docker.com/installation/)
- make

## :package: Dependencies

|NAME                          |URL|Version|
|-------------------------------|-----------------------------|--------|
|Gin|github.com/gin-gonic/gin|v1.9.0|
|Go playground validator|github.com/go-playground/validator/v10|v10.12.0|
|MySQL Go driver|github.com/go-sql-driver/mysql|v1.7.0|
|Logrus|github.com/sirupsen/logrus|v1.9.0|
|Testify|github.com/stretchr/testify|v1.8.2|

### Indirect Dependencies

> See [go.mod](go.mod)

## :vertical_traffic_light: Usages

### Remote Uses

The binary is built during deployment

- Possible(s) option(s) : `-config=[PATH]`

>By default if not specified, config path will be ./parameters/parameters.json

- Example :

````shell:
# run api
go run main.go

# run api with specific config
go run main.go -config=./parameters/parameters.json
````

### Start containers

````:shell
make
````

### Local Installation

 ````:shell
git clone https://github.com/arckadious/fizzbuzz.git

cd fizzbuzz

make install
````

> 'make install' : You will not have to download dependencies again while running 'make run' or 'make srv'.

### Run API Server (with air)

> [What is air ?](https://github.com/cosmtrek/air)

````:shell
make srv
````

> If you have issues with live reloading, try poll = true instead of false in [.air.toml](.air.toml) config file, or update docker to last version.

### Run API Server (without air)

````:shell
make run
````

### Run Tests

> MariaDB database need to be initialized and available ('make tests' use MySQL database from [api config file](./parameters/parameters.json)).

````:shell
make tests
````

>For more commands, see [Makefile](Makefile)

## :whale: Environment

|ENVIRONMENT|URL RANCHER|
|:--|:--|
|DEV|<fizbuzz-dev.example.com>|
|RECETTE|<https://fizbuzz-rct.example.com>|
|PROD|<https://fizbuzz.example.com>|

## :mag: Configuration

Check this config file [/parameters/parameters.json](/parameters/parameters.json)

:exclamation: Sections to change by environment
|Section|Description|
|--|--|
|env|api environment|
|database|mariadb database configuration|

## :page_with_curl: Logs

This application use 3 types of logs:

- Error logs in standard output and exit the application (also available from dozzle.localhost:8080)
- Applicatives logs stored in audit database
- Gin framework logs on localhost environment in gin.log file

### Applicative logs

Fizzbuzz API uses a logger middleware, which send requests and responses to a MySQL Database.

>'/swagger' api endpoint is excluded from applicative logs.

Using local database, you can see these logs at phpmyadmin.localhost:8080. Requests and responses are separated in two tables, and bound by a "COR_ID".

````sql
-- SQL Left Join example to use the COR_ID bind
SELECT APP_NAME, STATUS, SERVICE_ADDRESS, HOST, mr.MSG as REQUEST,
mp.MSG as RESPONSE, mr.COR_ID, mr.DT_CREATION as DT_CREATION_REQ,
mp.DT_CREATION as DT_CREATION_RESP
FROM MESSAGES_REQUEST as mr
LEFT JOIN MESSAGES_RESPONSE as mp ON mr.COR_ID = mp.COR_ID;
````

### Gin Framework Logs

Gin Framework log each request sent to fizzbuzz API.

> gin.log logs are available only for localhost environment purposes.

- Example :

````log
[GIN-debug] GET    /swagger/*filepath        --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (3 handlers)
[GIN-debug] HEAD   /swagger/*filepath        --> github.com/gin-gonic/gin.(*RouterGroup).createStaticHandler.func1 (3 handlers)
[GIN-debug] GET    /ping                     --> github.com/arckadious/fizzbuzz/server.(*Server).Handler.func4 (4 handlers)
[GIN-debug] POST   /v1/fizzbuzz              --> github.com/arckadious/fizzbuzz/server.(*Server).Handler.func5 (4 handlers)
[GIN-debug] GET    /v1/statistics            --> github.com/arckadious/fizzbuzz/server.(*Server).Handler.func6 (4 handlers)
[GIN] 2023/04/25 - 19:28:14 | 200 |        91.1µs |   192.168.208.5 | GET      "/ping"
[GIN] 2023/04/25 - 19:40:49 | 405 |        40.5µs |   192.168.208.5 | PUT  "/v1/fizzbuzz"
````

## :link: Tools

Local tools are available to make yourself comfortable :

- Local MySQL database + phpmyadmin web interface.
- Dozzle (log and memory monitoring for all containers).
- Swagger Web interface.

| ENV |URL|
|---|--|
|Swagger| <http://swagger.localhost:8080>
|Phpmyadmin| <http://phpmyadmin.localhost:8080>
|Dozzle (logs)| <http://dozzle.localhost:8080>
|API Go| <http://api.localhost:8080>

> For more details, see [nginx.conf](nginx.conf) and [Docker-compose.yml](docker-compose.yml)

## :trident: API Endpoints

### Main endpoint

- /v1/fizzbuzz

> Returns a list of strings with numbers from 1 to limit, where: all multiples specified are replaced by text.

<http://api.localhost:8080/v1/fizzbuzz>

### Most request used

- /v1/statistics

> Return the parameters corresponding to the most used request, as well as the number of hits for this request.

<http://api.localhost:8080/v1/statistics>

### Doc OpenAPI(swagger)

- /swagger

> Swagger Documentation available from the api. Also available from swagger.localhost:8080.

<http://api.localhost:8080/swagger>

### :green_heart: Health check

#### Ping

- /ping

> Ping to check if service is online

<http://api.localhost:8080/ping>

## :file_folder: Workspace

<!--generated on Windows with tree /F -->

````schema
fizzbuzz
    ├───.github
    │   └───workflows   >>> github workflows, for codacy and codecov
    ├───action
    │   └───fizz        >>> extract and validate data
    │
    ├───config          >>> init API configuration
    ├───constant
    ├───container       >>> init all class  
    │
    ├───database        >>> init db
    │
    ├───manager         >>> process data from action handlers
    ├───model
    ├───parameters      >>> JSON api configuration (default) directory 
    │
    ├───repository      >>> contains all functions which interact with database
    │
    ├───response        >>> JSON response templates functions
    │
    ├───server          >>> init Gin framework and specify endpoints, middlewares
    │
    ├───swaggerui       >>> API documentation (Cf. swagger.json)
    │
    ├───util      >>> contains various useful tools like MD5 Hash and http body extract
    │
    ├───tests     >>> contains 'make tests' results (with html coverage file)
    │
    └───validator
````

## :briefcase: Author

- [Pierre PERRIER-RIDET](https://fr.linkedin.com/in/pierre-perrier-ridet-3561b9139)
