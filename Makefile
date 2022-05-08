# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

GOLANG_CROSS_VERSION := v1.18.1

.PHONY: build docs

build:
	goreleaser release --snapshot --rm-dist

docs:
	@find ./docs -name '*.md' -delete
	go run github.com/atomix/cli/cmd/atomix gen-docs --markdown -o ./docs

reuse-tool: # @HELP install reuse if not present
	command -v reuse || python3 -m pip install reuse

license: reuse-tool # @HELP run license checks
	reuse lint
