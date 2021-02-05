A minimal Gemini server written in Go. Launch your Gemini capsule atop a Titan II rocket.

I did have Titan II down to 110 lines, including whitespace, but the addition of comments, 
logging and a couple of very necessary security fixes have subsequently bloated the code to 157 lines. 

# Installation

If you have `go` installed and a `GOPATH` configured, then clone the repo and run `go install`:

```sh
git clone https://gitlab.com/lostleonardo/titan2.git
cd titan2
go install
```

Or you can download a prebuilt Linux binary from the 'Releases' page of this project. 

# Configuration

Titan II is configured using command line arguments. Provided that you have `titan2` in your `PATH`,
you can run it like so:

```sh
titan2 -hostname my.site -dir /home/ll/gemini -crt /home/ll/certs/crt.pem -key /home/ll/certs/key.pem -port 1965
```

You can access the help via `titan2 -h`.

# Setting Up Your Own Gemini Server

Rather than provide a full tutorial on this page, I shall simply link to two of the best. 
[This](https://share.tube/videos/watch/4fe4e1f0-7896-4b8c-bfb8-2ff19c78d8e5) by
Chris Were and [this](https://share.tube/videos/watch/a44503e9-efdf-48ea-a30d-f5eec00214db) by Uoou are both excellent.

They demonstrate the process using [Agate](https://github.com/mbrubeck/agate), but it should be
straightfoward to adapt the steps. Most of what you need is
in this repository; simply change the configuration values to match your environment rather than mine.

# Thanks

Titan II takes its lead from [Go-Gemini](https://git.sr.ht/~yotam/go-gemini) and
[Melchior](https://github.com/praetoriansentry/melchior) and, of course, the 
[Gemini spec](https://gemini.circumlunar.space/docs/specification.html), which is a masterpiece of minimalism.

A week or so ago, I was discovering Mastadon and ActivityPub, and thinking about building an ActivityPub
server. Then, Gemini swept in and I was truly inspired. The minimalism, the design philosophy, the focus on security
and the DIY ethic are a breath of fresh air. Not to mention the aesthetics; Gemini capsules in Geminispace with
a protocol positioned somewhere between Mercury (Gopher) and Apollo (the Web). Great stuff.

In my little corner of cyberspace, there has been so much attention lavished on Gemini so quickly
that I can't now remember (it's only been a week!) exactly where I first heard or read about it, but it was probably
one of Chris Were, HexDSL or Uoou, who have produced some of the best videos on the subject.

# Specification, I don't need no stinkin' specification

Seriously though, the spec is obviously vital; none of the fun that people are having with servers
and browsers and reading and writing would be possible without the Gemini protocol.

However, minimal as it is, there are parts of the spec that I do not need and, as result, I have simply not implemented them. If I have overlooked anything essential, however, do please let me know.

## Input

I have no use for input at the moment and, to be honest, I am not entirely sure how it should work.

## Redirects

Likewise, I don't need a server that implements redirects; all I want to do is server static content from a single directory. I could add redirects, if required. Would parsing a text file that maps an old URL to a new URL fit the bill?

## Client Certs

The spec itself acknowledges that minimal implementations are unlikely to make use of this capability.

# Oh, just one more thing

Titan II successfully launched its first payload into Geminispace on 2021-02-04 at roughly 23:42 GMT.

My Gemini capsule is at: `gemini://gemini.lostleonardo.xyz`.
