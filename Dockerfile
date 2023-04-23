#build stage
ARG GO_VERSION
FROM golang:${GO_VERSION}-alpine3.17 AS builder
RUN apk add --no-cache make
WORKDIR /app
COPY . .
RUN make ondehoje

#final stage
FROM alpine:3.17.2
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/cmd/ondehoje /app/cmd/ondehoje
CMD ["/app/cmd/ondehoje/ondehoje"]
