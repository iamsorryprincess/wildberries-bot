lint:
	~/go/bin/golangci-lint run ./... -c .golangci.yaml

docker-build-api:
	docker build -f cmd/api/Dockerfile --build-arg DIR=api --build-arg PORT=8080 -t kekit/api .