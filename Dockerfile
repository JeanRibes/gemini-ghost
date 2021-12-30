FROM golang:1.17-alpine as builder

WORKDIR /go/src
RUN mkdir -p geminighost


COPY go.mod geminighost
COPY go.sum geminighost

# téléchargement des dépendances
RUN cd geminighost && go mod download

COPY *.go /go/src/geminighost/

# build Go
RUN cd geminighost && go build -o /main

#FROM alpine
#COPY --from=builder /main .

CMD ["/main"]