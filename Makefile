build:
	GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc-win32 go build -o ./dist/simcy.exe

clean:
	rm ./dist/simcy.exe

tests:
	go test ./... --cover

.PHONY: build clean tests