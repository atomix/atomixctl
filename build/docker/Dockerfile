# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: Apache-2.0

FROM alpine:3.13

RUN apk upgrade --update --no-cache && apk add bash bash-completion curl libc6-compat
RUN addgroup -S atomix && adduser -S -G atomix atomix

USER atomix
WORKDIR /home/atomix

COPY build/output/atomix_linux_amd64 /usr/local/bin/atomix

RUN mkdir -p /home/atomix/.atomix && \
    echo "" >> /home/atomix/.atomix/config.yaml && \
    cp /etc/profile /home/atomix/.bashrc && \
    atomix completion bash > /home/atomix/.atomix/bash_completion.sh && \
    echo "source /home/atomix/.atomix/bash_completion.sh" >> /home/atomix/.bashrc

ENTRYPOINT ["/bin/bash"]
