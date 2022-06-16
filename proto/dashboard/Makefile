.DEFAULT_GOAL := all
SERVICE := dashboard
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTO_PATH := $(subst /${SERVICE},,${CURDIR})
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

PROTOS := $(wildcard *.proto)
PROTO_PATHS = $(patsubst %.proto,$(SERVICE)/%.proto,$(PROTOS))

PROTO_PB = ${SERVICE}.pb
$(PROTO_PB): $(PROTOS)
	$(PROTOC) -I${PROTO_PATH} --include_imports --include_source_info \
    --descriptor_set_out=${SERVICE}.pb $(PROTO_PATHS)

OPENAPI := ../../openapi/${SERVICE}/openapi.yaml
$(OPENAPI): ${PROTOS} | ${PROTOC}
	protoc ${SERVICE}/${SERVICE}.proto -I${PROTO_PATH} --openapi_out=../../openapi/dashboard

OPENAPI_ENUMS := ../../openapi/dashboard/enums.yaml

# Go
GO_PROTO_PATH=../../go/proto
PROTO_BASES := $(patsubst %.proto,%,$(PROTOS))
GO_PATHS = $(GO_PROTO_PATH)/$(PROTO_BASE)/$(PROTO_BASE).pb.go
GO_GRPC_PATHS =   ../../go/proto/$(GRPC_BASE)/$(GRPC_BASE)_grpc.pb.go
GO_PROTOS := $(foreach PROTO_BASE,$(PROTO_BASES),$(GO_PATHS)) $(foreach GRPC_BASE,$(SERVICE),$(GO_GRPC_PATHS))

$(GO_PROTOS): $(PROTOS)
	$(PROTOC) -I${PROTO_PATH} --go_out=../../go --go_opt=module=platform/go \
			--go-grpc_out=../../go --go-grpc_opt=module=platform/go $(patsubst %.proto,$(SERVICE)/%.proto,$^)

# Typescript App
TS_PATH=../../app/api
TS_ENUM_PATH=$(addprefix $(TS_PATH)/,enums)
TS_FILE_NAMES=api.ts
TS_FILES=$(addprefix $(TS_PATH)/, $(TS_FILE_NAMES))
TS_ENUMS=$(addprefix $(TS_ENUM_PATH)/, $(TS_FILE_NAMES))

export TS_POST_PROCESS_FILE=../../app/node_modules/prettier/bin-prettier.js --write
$(TS_FILES) &:: $(OPENAPI)
	openapi-generator generate --remove-operation-id-prefix --enable-post-process-file -i $< -g typescript-axios -o $(TS_PATH)
$(TS_ENUMS) &:: $(OPENAPI_ENUMS)
	openapi-generator generate --skip-validate-spec --remove-operation-id-prefix --enable-post-process-file -i $< -g typescript-axios -o $(TS_ENUM_PATH)

clean:
	rm -f $(PROTO_PB)
	rm -f $(OPENAPI)
	rm -rf $(TS_PATH)
	rm -rf ${GO_PROTOS}

# For debugging
# $(info $$var is [${TS_FILES}])

all: $(PROTO_PB) $(GO_PROTOS) $(OPENAPI) $(TS_FILES) $(TS_ENUMS)
