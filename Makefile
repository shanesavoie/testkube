.PHONY: test cover 

BIN_DIR ?= $(HOME)/bin
GITHUB_TOKEN ?= "SET_ME"
USER ?= $(USER)
DATE ?= $(shell date -u --iso-8601=seconds)
COMMIT ?= $(shell git log -1 --pretty=format:"%h")

run-api-server: 
	APISERVER_PORT=8080 go run cmd/api-server/main.go

run-api-server-telepresence: 
	API_MONGO_DSN=mongodb://kubtest-mongodb:27017 POSTMANEXECUTOR_URI=http://kubtest-postman-executor:8082 APISERVER_PORT=8080 go run cmd/api-server/main.go


run-executor: 
	POSTMANEXECUTOR_PORT=8082 go run cmd/postman-executor/main.go

run-mongo-dev: 
	docker run -p 27017:27017 mongo


build: build-executor build-api-server build-kubtest-bin

# build done by vendoring to bypass private go repo problems
build-executor: 
	go build -o $(BIN_DIR)/postman-executor cmd/postman-executor/main.go

build-api-server:
	go build -o $(BIN_DIR)/api-server cmd/api-server/main.go 

build-kubtest-bin: 
	go build -ldflags="-s -w -X main.version=0.0.0-$(COMMIT) -X main.commit=$(COMMIT) -X main.date=$(DATE) -X main.builtBy=$(USER)" -o "$(BIN_DIR)/kubectl-kubtest" cmd/kubectl-kubtest/main.go


# build done by vendoring to bypass private go repo problems
docker-build-executor: 
	go mod vendor
	docker build --build-arg TOKEN=$(GITHUB_TOKEN) -t postman-executor -f build/postman-executor/Dockerfile .

docker-build-api-server:
	go mod vendor
	docker build --build-arg TOKEN=$(GITHUB_TOKEN) -t api-server -f build/api-server/Dockerfile .


install-swagger-codegen-mac: 
	brew install swagger-codegen

openapi-generate-model: openapi-generate-model-kubtest 

openapi-generate-model-kubtest:
	swagger-codegen generate -i api/v1/kubtest.yaml -l go -o tmp/api/kubtest
	mv tmp/api/kubtest/model_*.go pkg/api/kubtest
	rm -rf tmp
	find ./pkg/api/kubtest -type f -exec sed -i '' -e "s/package swagger/package kubtest/g" {} \;
	

test: 
	go test ./... -cover

cover: 
	@go test -mod=vendor -failfast -count=1 -v -tags test  -coverprofile=./testCoverage.txt ./... && go tool cover -html=./testCoverage.txt -o testCoverage.html && rm ./testCoverage.txt 
	open testCoverage.html

diagrams: 
	plantuml ./docs/puml/*.puml -o ../img/

version-bump: version-bump-patch

version-bump-patch:
	go run scripts/bump.go -kind patch

version-bump-minor:
	go run scripts/bump.go -kind minor

version-bump-major:
	go run scripts/bump.go -kind major

commands-reference: 
	go run cmd/kubectl-kubtest/main.go doc > ./docs/reference.md