# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

.PHONY: build docs

build:
	goreleaser release --snapshot --rm-dist

release:
	goreleaser release

docs:
	go run github.com/atomix/cli/cmd/atomix-generate-docs

reuse-tool: # @HELP install reuse if not present
	command -v reuse || python3 -m pip install reuse

license: reuse-tool # @HELP run license checks
	reuse lint
