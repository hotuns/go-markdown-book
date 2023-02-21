GO111MODULE=on

.PHONY: build
build: bindata markdown-book

.PHONY: docker-push
docker-push: package-all
	docker buildx build --platform linux/arm64,linux/amd64 -t hedongshu/markdown-book:latest . --push

.PHONY: docker-build
docker-build: package-all
	docker build -t hedongshu/markdown-book:dev -f ./Dockerfile.Develop .

.PHONY: bindata
bindata:
	go install github.com/go-bindata/go-bindata/v3/go-bindata@latest
	go generate ./...

.PHONY: markdown-book
markdown-book:
	go build $(RACE) -o bin/markdown-book ./

.PHONY: build-race
build-race: enable-race build

.PHONY: run
run: build
	./bin/markdown-book web --config ./config/config.yml

.PHONY: run-race
run-race: enable-race run

.PHONY: test
test:
	go test $(RACE) ./...

.PHONY: test-race
test-race: enable-race test

.PHONY: enable-race
enable-race:
	$(eval RACE = -race)

.PHONY: package
package: build
	bash ./package.sh

.PHONY: package-all
package-all: build
	bash ./package.sh -p 'linux darwin windows' -a 'amd64 arm64'

.PHONY: clean
clean:
	rm ./bin/markdown-book
