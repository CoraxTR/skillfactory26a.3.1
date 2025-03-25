FROM golang:1.23.5 AS compiling_stage
RUN mkdir -p /go/src/pipeline
WORKDIR /go/src/pipeline
COPY pipeline.go .
COPY go.mod .
RUN go install .

FROM alpine:latest
LABEL version = "1.1"
LABEL maintainer = "Boris Sapozhnikov <obsistway@gmail.com>"
WORKDIR /root/
COPY --from=compiling_stage /go/bin/pipeline .
ENTRYPOINT ./pipeline