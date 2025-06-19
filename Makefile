lint:
	~/go/bin/golangci-lint run ./... -c .golangci.yaml

docker-build-api:
	docker build -f cmd/api/Dockerfile --build-arg PORT=8080 -t wildberries-bot/api .

build-api:
	sh scripts/build.sh

# Dev
# Backend dev dir
backend_dev_dir = deploy/dev

dev-infrastructure-run:
	docker compose -f $(backend_dev_dir)/infrastructure.yaml -p wildberries-bot-dev-infrastructure up -d

dev-infrastructure-down:
	docker compose -f $(backend_dev_dir)/infrastructure.yaml -p wildberries-bot-dev-infrastructure down

dev-mysql-connect:
	mysql --login-path=wb-bot-dev -D testdb