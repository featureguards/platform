.DEFAULT_GOAL := all
SERVICE := private
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

# Go
# Can't do % because target has multiple % in it. Substitutions don't work.
GO_PROTO_PATH=../../go/proto
$(GO_PROTO_PATH)/${SERVICE}/internal_feature_toggle.pb.go &:: internal_feature_toggle.proto
	$(PROTOC) -I${PROTO_PATH} -I. --go_out=../../go --go_opt=module=platform/go  $<


PROTO_BASES := $(patsubst %.proto,%,$(PROTOS))
GO_PATHS = $(GO_PROTO_PATH)/$(SERVICE)/$(PROTO_BASE).pb.go
GO_PROTOS := $(foreach PROTO_BASE,$(PROTO_BASES),$(GO_PATHS))

clean:
	rm -f $(PROTO_PB)
	rm -rf ${GO_PROTOS}

# For debugging
# $(info $$var is [${GO_PROTOS}])

all: $(GO_PROTOS)
