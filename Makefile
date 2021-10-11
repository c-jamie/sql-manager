ENV = prod staging local

.PHONY: build clean deploy the sql-manager platform
SHELL=/bin/bash

ifeq (${ENV}, prod)
	include .prod.env
else 	
	ifeq (${ENV}, staging)
		include .staging.env
	else
		ifeq (${ENV}, dev)
			include .dev.env
		else
			ifeq (${ENV}, int)
				include .int.env
			else
				ifeq (${ENV}, local-dev)
					include .local.dev.env
				else 	
					ifeq (${ENV}, local-staging)
						include .local.staging.env
					endif
				endif
			endif
		endif
	endif
endif
export

check-staging:
ifeq ($(ENV),staging)
	@echo "STAGING ENV active"
else
	$(error == env must be staging==)	
endif

check-local-staging:
ifeq ($(ENV),local-staging)
	@echo "STAGING ENV active"
else
	$(error == env must be local-staging==)	
endif

check-prod:
ifeq ($(ENV),prod)
	@echo "PROD ENV active"
else
	$(error == env must be PROD==)	
endif

check-dev:
ifeq ($(ENV),dev)
else
	$(error == env must be DEV==)	
endif

check-int:
ifeq ($(ENV),int)
	@echo "LOCAL ENV active"
else
	$(error == env must be LOCAL DEV==)	
endif


#######################
# SERVER
#######################

# DEV

run-server: check-dev
	docker rm -f sqlm-server-dev || true && \
	docker run --name sqlm-server-dev -d -p $${SQLM_SER_DB_PORT}:5432 -e POSTGRES_USER=$${SQLM_SER_DB_USER} -e POSTGRES_DB=$${SQLM_SER_DB_NAME} -e POSTGRES_PASSWORD=$${SQLM_SER_DB_PW} postgres:latest && \
	sleep 2s && \
	go run ./cmd/sqlmserver/main.go

# INT

test-server: check-int
	docker rm -f sqlm-server-int || true && \
	docker run --name sqlm-server-int -d -p $${SQLM_SER_DB_PORT}:5432 -e POSTGRES_USER=$${SQLM_SER_DB_USER} -e POSTGRES_DB=$${SQLM_SER_DB_NAME} -e POSTGRES_PASSWORD=$${SQLM_SER_DB_PW} postgres:latest && \
	sleep 2s && \
	go test -count=1 -v ./serverlib/integration -test.run=${args}

test-server-ga:
	go test -count=1 ./serverlib/integration

# DB

test-db: check-dev
	docker rm -f test-db || true && \
	docker run --name test-db -d -p $${TEST_DB_PORT}:5432 -e POSTGRES_USER=$${TEST_DB_USER} -e POSTGRES_DB=$${TEST_DB_NAME} -e POSTGRES_PASSWORD=$${TEST_DB_PW} postgres:latest && \
	sleep 2s

#######################
# CLIENT
#######################

# BUILD

# DEV

run-client: check-dev
	@go run ./cmd/sqlmclient/main.go ${args}

# TEST

test-client:
	go test -count=1 -v ./clientlib/sql -run ${args}


#######################
# RELEASE
#######################

release:
	git tag -a v${tag} -m "release v${tag}" && git push --follow-tags

run-release-local:
	./dist/sql-manager-test_linux_amd64/sql-manager-auth