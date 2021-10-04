LDFLAGS=-w -X 'hawton.dev/hygieia/internal/version.BuildTime=$(shell TZ='UTC' date)' -X 'hawton.dev/hygieia/internal/version.GitCommit=$(shell git rev-parse --short HEAD)' -X 'hawton.dev/hygieia/internal/version.Version=$(shell git describe --tags --abbrev=0 HEAD || echo dev)' -X 'hawton.dev/hygieia/internal/version.GoVersion=$(shell go version | awk '{print $$3}')'
COMPILER=go build -ldflags="$(LDFLAGS)"

build:
	GOOS=windows GOARCH=amd64 $(COMPILER) -o build/hygieia-windows.exe
	GOOS=linux GOARCH=amd64 $(COMPILER) -o build/hygieia-linux-amd64
	GOOS=linux GOARCH=386 $(COMPILER) -o build/hygieia-linux-386
	GOOS=darwin GOARCH=amd64 $(COMPILER) -o build/hygieia-osx-amd64
	GOOS=darwin GOARCH=arm64 $(COMPILER) -o build/hygieia-osx-arm64
	GOOS=linux GOARCH=arm GOARM=6 $(COMPILER) -o build/hygieia-linux-arm
	GOOS=linux GOARCH=arm64 GOARM=7 $(COMPILER) -o build/hygieia-linux-arm64

clean:
	rm -rf build