# Fizzbuzz
 [![Codacy Badge](https://app.codacy.com/project/badge/Grade/d097d0142e6043a3936879cd0433a696)](https://app.codacy.com/gh/arckadious/fizzbuzz/dashboard?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade)  [![Go Report Card](https://goreportcard.com/badge/github.com/arckadious/fizzbuzz?refresh=1)](https://goreportcard.com/report/github.com/arckadious/fizzbuzz)  [![codecov](https://codecov.io/gh/arckadious/fizzbuzz/branch/master/graph/badge.svg)](https://codecov.io/gh/arckadious/fizzbuzz)  [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/arckadious/fizbuzz/LICENSE)  [![Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](https://go.dev/)



The original fizzbuzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by "fizz", all multiples of 5 by "buzz", and all multiples of 15 by "fizzbuzz".

Fizzbuzz is a web server. It has two REST API endpoints, one will return the list of numbers, and another one will show what the most frequent request has been.

______________
<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [:arrows_clockwise: Requirements](#arrows_clockwise-requirements)
      - [Local Requirements](#local-requirements)
- [:package: Dependencies](#package-dependencies)
- [:vertical_traffic_light: Usages](#vertical_traffic_light-usages)
    - [Remote Uses](#remote-uses)
    - [Start containers](#start-containers)
    - [Local Installation](#local-installation)
    - [Run Server (with air)](#run-server-with-air)
    - [Run Server (without air)](#run-server-without-air)
    - [Run Tests](#run-tests)
- [:whale: Environment](#whale-environment)
- [:mag: Configuration](#mag-configuration)
- [:page_with_curl: Logs](#page_with_curl-logs)
- [:link: Local development](#link-local-development)
- [:trident: API Endpoints <a id="endpoints"></a>](#trident-api-endpoints-a-idendpointsa)
  - [Main endpoint](#main-endpoint)
  - [Most request used](#most-request-used)
  - [Doc OpenAPI(swagger)](#doc-openapiswagger)
  - [:green_heart: Health check <a id="healthcheck"></a>](#green_heart-health-check-a-idhealthchecka)
    - [Ping](#ping)
- [:file_folder: Workspace](#file_folder-workspace)
- [:briefcase: Author](#briefcase-author)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->
<!-- MSYS_NO_PATHCONV=1 docker run --rm -v "$(pwd)":/app peterdavehello/npm-doctoc doctoc /app/README.md -->


## :arrows_clockwise: Requirements

Golang 1.19+

#####  Local Requirements
* [docker](https://docs.docker.com/installation/)
* [docker-compose](https://docs.docker.com/compose/install/)
* make

## :package: Dependencies

|NAME                          |URL|Version|
|-------------------------------|-----------------------------|--------|
|Gin|github.com/gin-gonic/gin|v1.9.0|
|Go playground validator|github.com/go-playground/validator/v10|v10.12.0|
|MySQL Go driver|github.com/go-sql-driver/mysql|v1.6.0|
|Logrus|github.com/sirupsen/logrus|v1.9.0|
 

## :vertical_traffic_light: Usages

#### Remote Uses

The binary is built during deployment

 - Possible option :
	  `-config=[PATH]`     => by default if not specified, config path will be ./parameters/parameters.json


 - Example :
````shell: 
   # run api
   $ go run main.go

   # run api with specific config
   $ go run main.go -config=./parameters/parameters.json
````
#### Start containers
````:shell 
$ make
  ````

#### Local Installation
 ````:shell 
$ git clone https://github.com/arckadious/fizzbuzz.git
$ cd fizzbuzz
$ make install
  ````


#### Run Server (with air)
> [What is air ?](https://github.com/cosmtrek/air)
````:shell 
$ make srv
  ````

#### Run Server (without air)
````:shell 
$ make run
  ````

#### Run Tests
````:shell 
$ make tests
  ````

>For more commands, see [Makefile](Makefile)

## :whale: Environment

|ENVIRONMENT|URL RANCHER|
|:--|:--|
|DEV|https://fizbuzz-dev.example.com|
|RECETTE|https://fizbuzz-rct.example.com|
|PROD|https://fizbuzz.example.com|

## :mag: Configuration

Check this config file [/parameters/parameters.json](/parameters/parameters.json)

:exclamation: Sections to change by environment
|Section|Description|
|--|--|
|env|api env|
|database|mariadb informations|

## :page_with_curl: Logs

This application use 2 types of logs:

 - Error log in standard output and exit the application
 - Applicatives logs in audit database, saved under database field APP_NAME=fizzbuzz

## :link: Local development

Local tools are available to make yourself comfortable :
- Local MySQL database + phpmyadmin web interface.
- Dozzle (log monitoring for all containers).
- Swagger Web interface.

| ENV |URL| 
|---|--|
|Swagger| http://swagger.localhost:8080
|Phpmyadmin| http://phpmyadmin.localhost:8080
|Dozzle (logs)| http://dozzle.localhost:8080
|API Go| http://api.localhost:8080

> For more details, see [nginx.conf](nginx.conf)

## :trident: API Endpoints <a id="endpoints"></a>


### Main endpoint
-  /v1/fizzbuzz

> Returns a list of strings with numbers from 1 to limit, where: all multiples specified are replaced by text.

http://api.localhost:8080/v1/fizzbuzz

### Most request used
-  /v1/statistics

> Return the parameters corresponding to the most used request, as well as the number of hits for this request.

http://api.localhost:8080/v1/statistics

### Doc OpenAPI(swagger)
- /swagger

> Swagger Documentation available from the api. Also available from swagger.localhost:8080.

http://api.localhost:8080/swagger

### :green_heart: Health check <a id="healthcheck"></a>
#### Ping
-  /ping

> Ping to check if service is online

http://api.localhost:8080/ping


## :file_folder: Workspace
<!--generated on Windows with tree /F -->

````
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
    ├───testing_results >>> contains 'make tests' results (with html coverage file)
    │
    └───validator
````
___________

## :briefcase: Author

* [Pierre PERRIER-RIDET](https://fr.linkedin.com/in/pierre-perrier-ridet-3561b9139)
