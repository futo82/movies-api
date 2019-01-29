FROM golang:1.11.5

RUN mkdir -p movies-api
WORKDIR /movies-api
ADD . /movies-api

RUN go get -u github.com/gin-gonic/gin
RUN go get -u github.com/aws/aws-sdk-go
RUN go get -u github.com/okta/okta-jwt-verifier-golang
RUN go build ./main.go

EXPOSE 8080

CMD ["./main"]