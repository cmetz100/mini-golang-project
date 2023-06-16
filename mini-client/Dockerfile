FROM golang:1.20.4
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o main .
CMD /app/main -url="http://$TASK_MANAGER_SERVICE_SERVICE_HOST:$TASK_MANAGER_SERVICE_SERVICE_PORT" -numIter=12 -numSec=150
