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
ENV URL "http://localhost:2368/ghost/api/v4/content/posts/"
ENV API_KEY "get from ghost admin -> new integration -> content api key"
CMD ["/main","-crt","/certs/crt.pem","-key","/certs/key.pem", "-hostname", "localhost", "-port", "1965"]
