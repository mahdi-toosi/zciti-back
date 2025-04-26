FROM golang:1.22.2-alpine3.19 as builder

WORKDIR /app

#COPY go.mod .

#RUN go mod tidy

COPY . .
COPY ./config/zcitiProd.toml ./config/zciti.toml

# Build the binary
RUN go build -mod=vendor -tags musl -o ./tmp/main.o ./cmd/main/main.go

## Use a minimal scratch image for the final stage
FROM alpine:3.19.1

## Copy the binary from the builder stage to the correct location
COPY --from=builder /app/tmp/main.o .
COPY --from=builder /app/config/zcitiProd.toml /app/config/zciti.toml

EXPOSE 8000
# Set the command to run the binary
#CMD ["go","run","cmd/main/main.go"]
CMD ["./main.o"]