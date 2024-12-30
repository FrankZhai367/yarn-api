FROM golang:1.23 as builder
ENV GOPROXY=https://goproxy.cn
WORKDIR /home/backend

COPY . .
RUN go build -o ./server main.go

FROM debian:stable-slim

RUN apt update && apt install -y libc6
# RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list && apt update && apt install -y libc6
WORKDIR /home/backend

COPY config.yaml .
COPY --from=builder /home/backend/server .
RUN echo 321 >> 1.txt
CMD tail -f 1.txt
# CMD ["./server"]

