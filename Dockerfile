FROM golang:1.11-alpine AS builder
RUN mkdir -p /go/fangs
WORKDIR /go/fangs

ENV GO111MODULE on
COPY . .

ENV CGO_ENABLED 0
RUN go build -a -installsuffix cgo -mod vendor -o /go/bin/fangs ./cmd/fangs

FROM scratch
COPY --from=builder /go/bin/fangs /go/bin/fangs
EXPOSE 8080:8080
ENTRYPOINT ["/go/bin/fangs"]