# List all subdirectories in cmd/ that (presumably) contain main.go
CMD_DIRS := $(wildcard cmd/*)

# Create output binary names:
#   e.g., cmd/foo  ->  bin/foo_linux_amd64
#          cmd/bar  ->  bin/bar_linux_amd64
#          cmd/foo  ->  bin/foo_darwin_arm64
# etc.
EXE_LINUX_AMD64   := $(CMD_DIRS:cmd/%=bin/%_linux_amd64)
EXE_DARWIN_ARM64  := $(CMD_DIRS:cmd/%=bin/%_darwin_arm64)

# Default target: build everything for both platforms.
all: build

## Build targets for both platforms
build: build-linux-amd64 build-darwin-arm64

# Build for linux/amd64
build-linux-amd64: $(EXE_LINUX_AMD64)

# Build for darwin/arm64
build-darwin-arm64: $(EXE_DARWIN_ARM64)

# Rule to build a single directory as a linux/amd64 binary
bin/%_linux_amd64: cmd/%
	GOOS=linux GOARCH=amd64 go build -o $@ ./$<

# Rule to build a single directory as a darwin/arm64 binary
bin/%_darwin_arm64: cmd/%
	GOOS=darwin GOARCH=arm64 go build -o $@ ./$<

# Clean out all generated binaries
clean:
	rm -rf bin/*

