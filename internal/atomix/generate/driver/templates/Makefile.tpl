.PHONY: build

build:
	goreleaser release --snapshot --rm-dist

release:
	goreleaser release

reuse-tool: # @HELP install reuse if not present
	command -v reuse || python3 -m pip install reuse

license: reuse-tool # @HELP run license checks
	reuse lint
