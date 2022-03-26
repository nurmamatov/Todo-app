FROM golang:1.17-alpine
RUN mkdir assignee-services
COPY . /assignee-services
WORKDIR /assignee-services
RUN go mod tidy
RUN go build -o main cmd/main.go
CMD ./main