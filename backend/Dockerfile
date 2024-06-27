FROM golang:1.22-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
#RUN apk update &&  apk add git
RUN go install github.com/air-verse/air@latest
#COPY . ./
#RUN go build -o main ./cmd/main.go

EXPOSE 8080
#CMD ["./main"]
CMD ["air", "-c", ".air.toml"]
