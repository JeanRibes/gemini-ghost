FROM --platform=$BUILDPLATFORM golang:1.17-alpine as builder

WORKDIR /go/src/gemini-ghost


COPY go.mod .
COPY go.sum .

# téléchargement des dépendances
RUN go mod download

COPY *.go ./
COPY ghost ./ghost/

# build Go
ARG TARGETOS TARGETARCH
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /main

#FROM alpine
FROM alpine:3.16
COPY --from=builder /main .
EXPOSE 1965
WORKDIR /
RUN adduser -DHu 1000 user
USER user

ADD index.tpl /
ADD static /static
ENV GHOST_URL "http://localhost:2368/ghost/api/v4/content"
ENV API_KEY "get from ghost admin -> new integration -> content api key"
CMD /main -crt /certs/crt.pem -key /certs/key.pem -ghost-url $GHOST_URL -ghost-key $API_KEY
