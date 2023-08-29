FROM amd64/alpine:latest

WORKDIR /work

ADD ./bin/linux-amd64-coin-manage /work/main

CMD ["./main"]

