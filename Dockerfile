FROM golang:alpine
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

RUN go build -o main ./bin/main.go


FROM alpine
WORKDIR /app

COPY --from=0 /app/main .
COPY --from=0 /app/config ./config

# Command to run the executable
CMD ["./main"]
