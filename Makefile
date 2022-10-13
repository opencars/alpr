.PHONY: default build clean
APPS        := server worker
BLDDIR      ?= bin
IMPORT_BASE := github.com/opencars/alpr

default: clean build

build: $(APPS)

$(BLDDIR)/%:
	go build -o $@ ./cmd/$*

$(APPS): %: $(BLDDIR)/%

lint:
	@revive -formatter stylish -config=revive.toml ./...

clean:
	@mkdir -p $(BLDDIR)
	@for app in $(APPS) ; do \
		rm -f $(BLDDIR)/$$app ; \
	done
