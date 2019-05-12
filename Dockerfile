FROM golang

WORKDIR $GOPATH/src/github.com/compujuckel/librefrontier
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN ls -la $GOPATH/bin/

CMD ["librefrontier"]
