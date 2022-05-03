# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

FROM docker.elastic.co/beats-dev/golang-crossbuild:1.18-darwin

RUN apt-get update && \
    apt-get install -y unzip git gcc

RUN mkdir /build && mkdir /resources

WORKDIR /build

COPY dist/atomix-build_linux_amd64_v1/atomix /usr/local/bin/atomix-build
COPY go.mod /resources/go.mod

ENTRYPOINT ["atomix-build"]
