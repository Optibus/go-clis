

build:
	go build -o ../bin/aws-secrets get-secrets.go

run:
	go run .

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/aws-secrets-linux-arm get-secrets.go
	GOOS=linux GOARCH=arm64 go build -o bin/aws-secrets-linux-arm64 get-secrets.go
	GOOS=freebsd GOARCH=386 go build -o bin/aws-secrets-freebsd-386 get-secrets.go
