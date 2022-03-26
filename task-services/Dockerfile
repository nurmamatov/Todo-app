FROM golang:1.17-alpine
RUN mkdir task-services
COPY . /task-services
WORKDIR /task-services
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main
