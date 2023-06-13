FROM golang:1.20.4
ENV DD_AGENT_HOST host.docker.internal 
RUN mkdir /app
RUN mkdir /app/logs
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD ["/app/main","/app/logs/mini-server-logs.log"]