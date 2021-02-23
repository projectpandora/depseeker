FROM golang:1.15-alpine3.13 AS build
RUN GO111MODULE=on go get -v github.com/projectpandora/depseeker/cmd/depseeker

FROM alpine:3.13
COPY --from=build /go/bin/depseeker /usr/local/bin/depseeker
ENTRYPOINT ["depseeker"]