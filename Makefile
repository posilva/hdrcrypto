hdrcrypto:
	go build ./cmd/hdrcrypto

.PHONY: run
run:
	go run ./cmd/hdrcrypto

.PHONY: docker-build
docker-build:
	docker build -t hdrcrypto . 

.PHONY: docker-run
docker-run: docker-build
	docker run -it hdrcrypto 

