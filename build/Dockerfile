FROM alpine:3.8

RUN apk upgrade --update --no-cache && apk add bash bash-completion
RUN addgroup -S atomix && adduser -S -G atomix atomix

USER atomix
WORKDIR /home/atomix

ADD build/_output/bin/atomix /usr/local/bin/atomix

RUN mkdir -p /home/atomix/.atomix && \
    echo "" >> /home/atomix/.atomix/config.yaml && \
    cp /etc/profile /home/atomix/.bashrc && \
    atomix completion bash > /home/atomix/.atomix/bash_completion.sh && \
    echo "source /home/atomix/.atomix/bash_completion.sh" >> /home/atomix/.bashrc

ENTRYPOINT ["/bin/bash"]
