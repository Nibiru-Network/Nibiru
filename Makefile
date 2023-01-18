# This Makefile is meant to be used by people that do not usually work
# with Go source code. If you know what GOPATH is then you probably
# don't need to bother with make.

.PHONY: nbn android ios nbn-cross evm all test clean
.PHONY: nbn-linux nbn-linux-386 nbn-linux-amd64 nbn-linux-mips64 nbn-linux-mips64le
.PHONY: nbn-linux-arm nbn-linux-arm-5 nbn-linux-arm-6 nbn-linux-arm-7 nbn-linux-arm64
.PHONY: nbn-darwin nbn-darwin-386 nbn-darwin-amd64
.PHONY: nbn-windows nbn-windows-386 nbn-windows-amd64

GOBIN = ./build/bin
GO ?= latest
GORUN = env GO111MODULE=on go run

nbn:
	$(GORUN) build/ci.go install ./cmd/nbn
	@echo "Done building."
	@echo "Run \"$(GOBIN)/nbn\" to launch nbn."

all:
	$(GORUN) build/ci.go install

android:
	$(GORUN) build/ci.go aar --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/nbn.aar\" to use the library."
	@echo "Import \"$(GOBIN)/nbn-sources.jar\" to add javadocs"
	@echo "For more info see https://stackoverflow.com/questions/20994336/android-studio-how-to-attach-javadoc"

ios:
	$(GORUN) build/ci.go xcode --local
	@echo "Done building."
	@echo "Import \"$(GOBIN)/nbn.framework\" to use the library."

test: all
	$(GORUN) build/ci.go test

lint: ## Run linters.
	$(GORUN) build/ci.go lint

clean:
	env GO111MODULE=on go clean -cache
	rm -fr build/_workspace/pkg/ $(GOBIN)/*

# The devtools target installs tools required for 'go generate'.
# You need to put $GOBIN (or $GOPATH/bin) in your PATH to use 'go generate'.

devtools:
	env GOBIN= go get -u golang.org/x/tools/cmd/stringer
	env GOBIN= go get -u github.com/kevinburke/go-bindata/go-bindata
	env GOBIN= go get -u github.com/fjl/gencodec
	env GOBIN= go get -u github.com/golang/protobuf/protoc-gen-go
	env GOBIN= go install ./cmd/abigen
	@type "npm" 2> /dev/null || echo 'Please install node.js and npm'
	@type "solc" 2> /dev/null || echo 'Please install solc'
	@type "protoc" 2> /dev/null || echo 'Please install protoc'

# Cross Compilation Targets (xgo)

nbn-cross: nbn-linux nbn-darwin nbn-windows nbn-android nbn-ios
	@echo "Full cross compilation done:"
	@ls -ld $(GOBIN)/nbn-*

nbn-linux: nbn-linux-386 nbn-linux-amd64 nbn-linux-arm nbn-linux-mips64 nbn-linux-mips64le
	@echo "Linux cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-*

nbn-linux-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/386 -v ./cmd/nbn
	@echo "Linux 386 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep 386

nbn-linux-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/amd64 -v ./cmd/nbn
	@echo "Linux amd64 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep amd64

nbn-linux-arm: nbn-linux-arm-5 nbn-linux-arm-6 nbn-linux-arm-7 nbn-linux-arm64
	@echo "Linux ARM cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep arm

nbn-linux-arm-5:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-5 -v ./cmd/nbn
	@echo "Linux ARMv5 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep arm-5

nbn-linux-arm-6:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-6 -v ./cmd/nbn
	@echo "Linux ARMv6 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep arm-6

nbn-linux-arm-7:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm-7 -v ./cmd/nbn
	@echo "Linux ARMv7 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep arm-7

nbn-linux-arm64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/arm64 -v ./cmd/nbn
	@echo "Linux ARM64 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep arm64

nbn-linux-mips:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips --ldflags '-extldflags "-static"' -v ./cmd/nbn
	@echo "Linux MIPS cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep mips

nbn-linux-mipsle:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mipsle --ldflags '-extldflags "-static"' -v ./cmd/nbn
	@echo "Linux MIPSle cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep mipsle

nbn-linux-mips64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64 --ldflags '-extldflags "-static"' -v ./cmd/nbn
	@echo "Linux MIPS64 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep mips64

nbn-linux-mips64le:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=linux/mips64le --ldflags '-extldflags "-static"' -v ./cmd/nbn
	@echo "Linux MIPS64le cross compilation done:"
	@ls -ld $(GOBIN)/nbn-linux-* | grep mips64le

nbn-darwin: nbn-darwin-386 nbn-darwin-amd64
	@echo "Darwin cross compilation done:"
	@ls -ld $(GOBIN)/nbn-darwin-*

nbn-darwin-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/386 -v ./cmd/nbn
	@echo "Darwin 386 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-darwin-* | grep 386

nbn-darwin-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=darwin/amd64 -v ./cmd/nbn
	@echo "Darwin amd64 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-darwin-* | grep amd64

nbn-windows: nbn-windows-386 nbn-windows-amd64
	@echo "Windows cross compilation done:"
	@ls -ld $(GOBIN)/nbn-windows-*

nbn-windows-386:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/386 -v ./cmd/nbn
	@echo "Windows 386 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-windows-* | grep 386

nbn-windows-amd64:
	$(GORUN) build/ci.go xgo -- --go=$(GO) --targets=windows/amd64 -v ./cmd/nbn
	@echo "Windows amd64 cross compilation done:"
	@ls -ld $(GOBIN)/nbn-windows-* | grep amd64
