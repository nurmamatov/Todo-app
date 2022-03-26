FROM golang:1.17-alpine
RUN mkdir email_service
COPY . /email_service
WORKDIR /email_service
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main