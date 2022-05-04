# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

FROM alpine:3.13

RUN apk upgrade --update --no-cache && apk add bash bash-completion curl libc6-compat
RUN addgroup -S atomix && adduser -S -G atomix atomix

USER atomix
WORKDIR /home/atomix

COPY atomix /usr/local/bin/atomix

ENTRYPOINT ["/bin/bash"]
