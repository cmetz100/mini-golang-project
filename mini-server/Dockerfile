FROM golang:1.20.4
ENV DD_AGENT_HOST host.docker.internal 
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]