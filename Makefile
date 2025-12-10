# Makefile for md-serve

BINARY=md-serve
PLATFORMS=darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# Setup linker flags for static linking
LDFLAGS=-ldflags "-extldflags -static"

# Build for current platform
build:
	go build -o ${BINARY} ./cmd/md-serve

# Build statically linked binary for current platform
build-static:
	CGO_ENABLED=0 go build ${LDFLAGS} -o ${BINARY} ./cmd/md-serve

# Build for all platforms
build-all: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(word 1, $(subst /, ,$@)) GOARCH=$(word 2, $(subst /, ,$@)) CGO_ENABLED=0 go build ${LDFLAGS} -o ${BINARY}-$(word 1, $(subst /, ,$@))-$(word 2, $(subst /, ,$@)) ./cmd/md-serve

# Install dependencies
deps:
	go mod tidy

# Run the server
run:
	go run ./cmd/md-serve

# Clean build artifacts
clean:
	rm -f ${BINARY}*