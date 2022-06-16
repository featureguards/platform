.DEFAULT_GOAL := all
SUBDIRS = auth shared toggles
ALL_SUBDIRS = $(SUBDIRS) private dashboard
ALL_SUBDIRS_CLEAN = $(addsuffix .clean,$(ALL_SUBDIRS))
JS_PROTO_PATH = ../../featureguards-js/shared/proto
PROTO_WILDCARD = $(wildcard $(SUBDIR)/*.proto)
PROTOS = $(foreach SUBDIR,$(SUBDIRS),$(PROTO_WILDCARD))

$(ALL_SUBDIRS): 
	$(MAKE) -C $@

$(ALL_SUBDIRS_CLEAN):
	$(MAKE) -C $(basename $@) clean

JS_FILES = $(patsubst %.proto,%.ts,$(PROTOS))
JS_PATHS = $(JS_PROTO_PATH)/$(JS_FILE)
JS_PROTOS = $(foreach JS_FILE,$(JS_FILES),$(JS_PATHS))

# $(info $$var is [${ALL_SUBDIRS}])
# $(info $$var is [${JS_PROTOS}])

$(JS_PROTOS):
	npx protoc --ts_out $(JS_PROTO_PATH)  --ts_opt optimize_code_size --proto_path . $(PROTOS)  

js: $(JS_PROTOS)

.PHONY: subdirs $(ALL_SUBDIRS)

subdirs: $(ALL_SUBDIRS)

$(ALL_SUBDIRS):
	$(MAKE) -C $@

all: subdirs js

clean: $(ALL_SUBDIRS_CLEAN)
	rm -rf $(JS_PROTO_PATH)/*
