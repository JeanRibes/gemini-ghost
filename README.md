```

```

A Gemini dynamic server that serves content from the SQLite database of a [Ghost](https://ghost.org/) installation. It
scans the Ghost database, and serves (published) posts by converting the HTML content into Gemini

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
docker run -d -p 1965:1965 -v $(pwd)/certs:/certs -v $(pwd)/ghost.sqlite3:/ghost.db --name gemini-ghost ghcr.io/jeanribes/gemini-ghost:master 
```

# Configuration

```shell
Usage of ./gemini-ghost:
  -crt string
        cert filename (default "./certs/crt.pem")
  -dbfile string
        SQLITE file to use (default "ghost.db")
  -hostname string
        hostname (default "localhost")
  -key string
        key filename (default "./certs/key.pem")
  -port int
        port number (default 1965)
```
