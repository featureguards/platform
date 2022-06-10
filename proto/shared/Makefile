.DEFAULT_GOAL := all
SERVICE := shared
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

# Go
# Can't do % because target has multiple % in it. Substitutions don't work.
GO_ROOT_PATH=${ROOT_PATH}/../featureguards-go
GO_PROTO_PATH=${GO_ROOT_PATH}/proto
$(GO_PROTO_PATH)/feature_toggle/feature_toggle.pb.go &:: feature_toggle.proto
	$(PROTOC) -I${PROTO_PATH} -I. --go_out=${GO_ROOT_PATH} --go_opt=module=github.com/featureguards/featureguards-go/v2 $<

$(GO_PROTO_PATH)/user/user.pb.go &:: user.proto
	$(PROTOC) -I${PROTO_PATH} -I. --go_out=${GO_ROOT_PATH} --go_opt=module=github.com/featureguards/featureguards-go/v2 $<


PROTO_BASES := $(patsubst %.proto,%,$(PROTOS))
GO_PATHS = $(GO_PROTO_PATH)/$(PROTO_BASE)/$(PROTO_BASE).pb.go
GO_PROTOS := $(foreach PROTO_BASE,$(PROTO_BASES),$(GO_PATHS))

clean:
	rm -f $(GO_PROTOS)

all: $(GO_PROTOS)
