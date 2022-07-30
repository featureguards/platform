.DEFAULT_GOAL := all
SERVICE := auth
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
PROTO_PATHS = $(patsubst %.proto,$(SERVICE)/%.proto,$(PROTOS))
PROTO_PB = ../../envoy/${SERVICE}.pb
$(PROTO_PB): $(PROTOS)
	$(PROTOC) -I${PROTO_PATH} --include_imports --include_source_info \
    --descriptor_set_out=$(PROTO_PB) $(PROTO_PATHS)

# Go
GO_ROOT_PATH=${ROOT_PATH}/../featureguards-go
GO_PROTO_PATH=${GO_ROOT_PATH}/proto
PROTO_BASES := $(patsubst %.proto,%,$(PROTOS))
GO_PATHS = $(GO_PROTO_PATH)/$(PROTO_BASE)/$(PROTO_BASE).pb.go
GO_GRPC_PATHS = ${GO_PROTO_PATH}/$(GRPC_BASE)/$(GRPC_BASE)_grpc.pb.go
GO_PROTOS := $(foreach PROTO_BASE,$(PROTO_BASES),$(GO_PATHS)) $(foreach GRPC_BASE,$(SERVICE),$(GO_GRPC_PATHS))

$(GO_PROTOS): $(PROTOS)
	$(PROTOC) -I${PROTO_PATH} --go_out=${GO_ROOT_PATH} --go_opt=module=github.com/featureguards/featureguards-go/v2 \
			--go-grpc_out=${GO_ROOT_PATH} --go-grpc_opt=module=github.com/featureguards/featureguards-go/v2 $(patsubst %.proto,$(SERVICE)/%.proto,$^)


# Python
PY_ROOT_PATH=${ROOT_PATH}/../featureguards-python/featureguards
PY_PROTO_PATH=${PY_ROOT_PATH}/proto
PY_PATHS = $(PY_PROTO_PATH)/$(PROTO_BASE)/$(PROTO_BASE)_pb2.py
PY_GRPC_PATHS = ${PY_PROTO_PATH}/$(GRPC_BASE)/$(GRPC_BASE)_pb2_grpc.py
PY_PROTOS := $(foreach PROTO_BASE,$(PROTO_BASES),$(PY_PATHS)) $(foreach GRPC_BASE,$(SERVICE),$(PY_GRPC_PATHS))

$(PY_PROTOS): $(PROTOS)
	python3 -m grpc_tools.protoc -I${PROTO_PATH} --mypy_out=${PY_PROTO_PATH} --python_out=${PY_PROTO_PATH} --mypy_grpc_out=${PY_PROTO_PATH} --grpc_python_out=${PY_PROTO_PATH} $(patsubst %.proto,$(SERVICE)/%.proto,$^)


clean:
	rm -f $(PROTO_PB)
	rm -f ${GO_PROTOS}
	rm -f ${PY_PROTOS}

# For debugging
# $(info $$var is [${TS_FILES}])

all: $(PROTO_PB) $(GO_PROTOS) $(PY_PROTOS)