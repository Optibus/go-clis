
build:
	go build -o ../bin/check-js-deps main.go

run:
	go run .

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o bin/check-js-deps-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/check-js-deps-linux-arm64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/check-js-deps-freebsd-386 main.go
