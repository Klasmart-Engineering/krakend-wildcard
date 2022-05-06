all: clean build

build:
	go build -buildmode=plugin -o <plugin-name>.so .

clean:
	go clean