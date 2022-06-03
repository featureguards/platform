.DEFAULT_GOAL := all
SERVICE := toggles
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTO_PATH := $(subst /${SERVICE},,${CURDIR})
ROOT_PATH := $(subst /proto/${SERVICE},,${CURDIR})
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
PROTO_PB = ${SERVICE}.pb
$(PROTO_PB): $(PROTOS)
	$(PROTOC) -I${PROTO_PATH} -I. --include_imports --include_source_info \
    --descriptor_set_out=${SERVICE}.pb $^

OPENAPI := ../../openapi/${SERVICE}/openapi.yaml
$(OPENAPI): ${PROTOS} | ${PROTOC}
	protoc ${SERVICE}.proto -I${PROTO_PATH} -I. --openapi_out=../../openapi/${SERVICE}

# Go
# Can't do % because target has multiple % in it. Substitutions don't work.
GO_ROOT_PATH=${ROOT_PATH}/../client-go
GO_PROTO_PATH=${GO_ROOT_PATH}/proto
$(GO_PROTO_PATH)/${SERVICE}/${SERVICE}.pb.go $(GO_PROTO_PATH)/${SERVICE}/${SERVICE}_grpc.pb.go &:: ${SERVICE}.proto
	$(PROTOC) -I${PROTO_PATH} -I. --go_out=${GO_ROOT_PATH} --go_opt=module=github.com/featureguards/client-go \
			--go-grpc_out=${GO_ROOT_PATH} --go-grpc_opt=module=github.com/featureguards/client-go $<

PROTO_BASES := $(patsubst %.proto,%,$(PROTOS))
GO_PATHS = $(GO_PROTO_PATH)/$(PROTO_BASE)/$(PROTO_BASE).pb.go
GO_GRPC_PATHS = ${GO_PROTO_PATH}/$(GRPC_BASE)/$(GRPC_BASE)_grpc.pb.go
GO_PROTOS := $(foreach PROTO_BASE,$(PROTO_BASES),$(GO_PATHS)) $(foreach GRPC_BASE,$(SERVICE),$(GO_GRPC_PATHS))

# Typescript Api
TS_PATH=../../sdk/typescript/${SERVICE}
TS_FILE_NAMES=api.ts
TS_FILES=$(addprefix $(TS_PATH)/, $(TS_FILE_NAMES))

export TS_POST_PROCESS_FILE=../../app/node_modules/prettier/bin-prettier.js --write
$(TS_FILES) &:: $(OPENAPI)
	openapi-generator generate --remove-operation-id-prefix --enable-post-process-file -i $< -g typescript-axios -o $(TS_PATH)

clean:
	rm -f $(PROTO_PB)
	rm -f $(OPENAPI)
	rm -rf $(TS_PATH)
	rm -f ${GO_PROTOS}

# For debugging
# $(info $$var is [${TS_FILES}])

all: $(PROTO_PB) $(GO_PROTOS) $(OPENAPI) $(TS_FILES)