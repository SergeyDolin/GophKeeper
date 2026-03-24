build:
	GOOS=linux GOARCH=amd64 go build -o bin/gophkeeper-linux
	GOOS=windows GOARCH=amd64 go build -o bin/gophkeeper.exe
	GOOS=darwin GOARCH=amd64 go build -o bin/gophkeeper-mac