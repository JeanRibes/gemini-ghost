A Gemini dynamic server that serves content from the Content API of a [Ghost](https://ghost.org/) installation. It
fetches all the posts, and serves them by converting the HTML content into Gemini

Based on the work of [lostleonardo's titan2](https://gitlab.com/lostleonardo/titan2.git)
, [gemini://gemini.lostleonardo.xyz](gemini://gemini.lostleonardo.xyz)

HTML-to-GMI conversion is done with [github.com/LukeEmmet/html2gemini](https://github.com/LukeEmmet/html2gemini)

# Installation

If you have `go` installed and a `GOPATH` configured, then clone the repo and run `go install`:

```sh
go install github.com/JeanRibes/gemini-ghost@master
```

## With Docker

```shell
mkdir certs
cd certs
openssl req -new -x509 -days 3650 -nodes -out crt.pem -keyout key.pem
cd ..
docker run -d -p 1965:1965 -e URL="http://localhost:2368/ghost/api/v4/content" -e API_KEY="<your ghost content api key here>" -v $(pwd)/certs:/certs --name gemini-ghost ghcr.io/jeanribes/gemini-ghost:master 
```

# Configuration

```shell
Usage of ./gemini-ghost:
  -crt string
        cert filename (default "./certs/crt.pem")
  -ghost-key string
        Ghost Content API Key (default "a513a3dc949855fb654a545bd7")
  -ghost-url string
        Ghost Content API Url (default "http://localhost:2368/ghost/api/v4/content")
  -hostname string
        hostname (default "localhost")
  -key string
        key filename (default "./certs/key.pem")
  -port int
        port number (default 1965)

```

# TODO

utiliser la navigation
http://localhost:2368/ghost/api/v4/content/settings/?key=a513a3dc949855fb654a545bd7&limit=3&page=1