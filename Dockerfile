FROM golang AS build

ARG GIT_DESC=undefined

WORKDIR /go/src/github.com/Snawoot/httptrap
COPY . .
RUN CGO_ENABLED=0 go build -a -tags netgo -ldflags '-s -w -extldflags "-static" -X main.version='"$GIT_DESC"

FROM scratch
COPY --from=build /go/src/github.com/Snawoot/httptrap/httptrap /
USER 9999:9999
EXPOSE 8008/tcp
ENTRYPOINT ["/httptrap", "-bind", "0.0.0.0:8008"]
