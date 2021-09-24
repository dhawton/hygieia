if [[ ! -d "build" ]]; then
    mkdir build
fi

GOOS=windows GOARCH=amd64 go build -o build/hygieia-windows.exe
GOOS=linux GOARCH=amd64 go build -o build/hygieia-linux-amd64
GOOS=linux GOARCH=386 go build -o build/hygieia-linux-386
GOOS=darwin GOARCH=amd64 go build -o build/hygieia-osx-amd64
GOOS=darwin GOARCH=arm64 go build -o build/hygieia-osx-arm64
GOOS=linux GOARCH=arm GOARM=6 go build -o build/hygieia-linux-arm
GOOS=linux GOARCH=arm64 GOARM=7 go build -o build/hygieia-linux-arm64