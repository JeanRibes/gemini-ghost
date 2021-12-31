FROM golang:1.17-alpine as builder

WORKDIR /go/src/gemini-ghost


COPY go.mod .
COPY go.sum .

# téléchargement des dépendances
RUN go mod download

COPY *.go .

# build Go
RUN go build -o /main

#FROM alpine
#COPY --from=builder /main .
EXPOSE 1965
WORKDIR /
RUN adduser -DHu 1000 grissom
USER grissom
ADD index.tpl /
CMD ["/main","-crt","/certs/crt.pem","-key","/certs/key.pem", "-dbfile", "/ghost.db", "-hostname", "localhost", "-port", "1965"]
