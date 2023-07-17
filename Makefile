before-push:
	go mod tidy &&\
	gofumpt -l -w . &&\
	go build ./...&&\
	golangci-lint run ./... &&\
	go test -v ./integration_tests/...


scrap:
	go run . scrap -o output
ingest:
	go run . ingest -i ./scrapper/output --email admin@mail.com --password "&dm1Npa$$"