build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ./bot ./main.go
	@echo "Build Success"