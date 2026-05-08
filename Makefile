test:
	go test ./microservices/greenSeeds/internal/transport/tests -coverpkg=./microservices/greenSeeds/internal/transport/ -coverprofile=coverage.out

cover:
	go tool cover -func=coverage.out

cover-html:
	go tool cover -html=coverage.out