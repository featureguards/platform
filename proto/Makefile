.DEFAULT_GOAL := all

PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go

PROTOC := $(shell which protoc)
# If protoc isn't on the path, set it to a target that's never up to date, so
# the install command always runs.
ifeq ($(PROTOC),)
    PROTOC = must-rebuild
endif

# Figure out which machine we're running on.
UNAME := $(shell uname)

$(PROTOC):
# Run the right installation command for the operating system.
ifeq ($(UNAME), Darwin)
	brew install protobuf
endif
ifeq ($(UNAME), Linux)
	sudo apt-get install protobuf-compiler
endif

# If $GOPATH/bin/protoc-gen-go does not exist, we'll run this command to install
# it.
$(PROTOC_GEN_GO):
	go get -u github.com/golang/protobuf/protoc-gen-go

PROTOS := $(shell find . -name '*.proto')
proto.pb: $(PROTOS) | $(PROTOC)
	$(PROTOC) -I${GOOGLEAPIS_DIR} -I. --include_imports --include_source_info \
    --descriptor_set_out=proto.pb *.proto

GO_PROTOS := $(patsubst %.proto,../go/proto/%,$(PROTOS))
$(GO_PROTOS): $(shell find . -name '*.proto') | $(PROTOC)
	$(PROTOC) --go_out=../go --go_opt=module=stackv2/go \
			--go-grpc_out=../go --go-grpc_opt=module=stackv2/go *.proto

clean:
	rm *.pb
	rm -rf ../go/proto/*

all: proto.pb $(GO_PROTOS)
