ifeq ($(OS),Windows_NT)
	HTML_TEST_OPEN = start "" "./tests/cover.html"
else 
	HTML_TEST_OPEN = open ./tests/cover.html
endif
.PHONY: all clean down build start start-rp start-db start-logs restart stop stop-rp stop-db stop-logs kill rm in install update srv run bash sh tests test

MODULE_NAME = github.com/arckadious/fizzbuzz
BUILD_FILENAME = fizzbuzz
DOCKER_COMPOSE_BIN = docker compose
DOCKER_COMPOSE = $(DOCKER_COMPOSE_BIN)

#Service Names on docker-compose
APP = api
REVERSEPROXY = nginx
DB = mariadb
LOGS = dozzle
PHPMYADMIN = phpmyadmin

DK_EXEC = docker exec -ti $$($(DOCKER_COMPOSE) ps -qa $(APP))

all: 
	@$(DOCKER_COMPOSE) up -d

clean:
	@$(DOCKER_COMPOSE) down --volumes --rmi local

down:
	@$(DOCKER_COMPOSE) down

build:
	@$(DOCKER_COMPOSE) build --no-cache

start:
	@$(DOCKER_COMPOSE) up --force-recreate -d $(APP)

start-rp:
	@$(DOCKER_COMPOSE) up -d $(REVERSEPROXY)

start-db: 
	@$(DOCKER_COMPOSE) up -d $(DB)
	@$(DOCKER_COMPOSE) up -d $(PHPMYADMIN)

start-logs: 
	@$(DOCKER_COMPOSE) up -d $(LOGS)

restart: stop start

stop:
	@$(DOCKER_COMPOSE) stop $(APP) || true

# stop reverse proxy if running and rm container
stop-rp:
	@docker stop $$($(DOCKER_COMPOSE) ps -qa $(REVERSEPROXY)) || true && docker rm --force $$($(DOCKER_COMPOSE) ps -qa $(REVERSEPROXY)) || true

stop-db:
	@docker stop $$($(DOCKER_COMPOSE) ps -qa $(DB)) || true && docker rm --force $$($(DOCKER_COMPOSE) ps -qa $(DB)) || true
	@docker stop $$($(DOCKER_COMPOSE) ps -qa $(PHPMYADMIN)) || true && docker rm --force $$($(DOCKER_COMPOSE) ps -qa $(PHPMYADMIN)) || true

stop-logs:
	@docker stop $$($(DOCKER_COMPOSE) ps -qa $(LOGS)) || true && docker rm --force $$($(DOCKER_COMPOSE) ps -qa $(LOGS)) || true

kill:
	@$(DOCKER_COMPOSE) kill || true

rm: stop
	@$(DOCKER_COMPOSE) rm --force $(APP) || true
	
in: install
install: all
	@$(DK_EXEC) bash -c "go mod tidy && go mod download && go mod vendor"

update:
	@$(DK_EXEC) bash -c "go get -u && go mod tidy && go mod download && go mod vendor"

srv:
	@$(DK_EXEC) bash -c "script -aqf /var/log/messages.log -c \"BUILD_FILENAME=${BUILD_FILENAME} air -c .air.toml\""

run:
	@$(DK_EXEC) bash -c "script -aqf /var/log/messages.log -c 'go build -o ${BUILD_FILENAME} . && ./${BUILD_FILENAME} -config=/go/config/config.json'"

bash:
	@$(DK_EXEC) bash

sh:
	@$(DK_EXEC) sh

tests:
	@$(DK_EXEC) bash -c "go test -v -coverprofile tests/cover.out ./... && go tool cover -html tests/cover.out -o tests/cover.html"
	@$(HTML_TEST_OPEN)
	
test: # example : ARGS=repository make test -> run package tests in repository folder
	@$(DK_EXEC) bash -c "go test -timeout 30s ${MODULE_NAME}/${ARGS}"