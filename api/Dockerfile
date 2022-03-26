FROM golang:1.17-alpine
RUN mkdir api
COPY . /api
WORKDIR /api
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main