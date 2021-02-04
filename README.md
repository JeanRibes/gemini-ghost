A minimal Gemini server written in Go. Launch your Gemini capsule atop a Titan II rocket.

I did have the code down to 110 lines including whitespace, but the addition of comments, logging and
a couple of very necessary security fixes have bloated the code to 155 lines. 

If you see any obvious flaws, please point them out.

# Installation

If you have `go` installed and a `GOPATH` configured, then clone the repo and run `go install`:

```sh
git clone https://gitlab.com/lostleonardo/titan2.git
cd titan2
go install
```

Or you can download a prebuilt binary from the releases page of this project. 

# Configuration

Titan II is configured using command line arguments. Provided that you have `titan2` in your `PATH`,
you can run it like so:

```sh
titan2 -hostname my.site -dir ./certs/gemini -crt ./certs/crt.pem -key ./crt.pem -port 1965
```

You can access the help to remind yourself of those options and what they mean via `titan2 -h`.

# Setting Up Your Own Gemini Server

Rather than providing a full tutorial on these pages, I shall simply link to two of the best. 
[This](https://share.tube/videos/watch/4fe4e1f0-7896-4b8c-bfb8-2ff19c78d8e5) by
Chris Were and [this](https://share.tube/videos/watch/a44503e9-efdf-48ea-a30d-f5eec00214db) by Uoou are both excellent.

They demonstrate the process using [Agate](https://github.com/mbrubeck/agate), but it should be
straightfoward to launch your Gemini capsule atop Titan II.

# Special Thanks

Titan II takes inspiration from both [Go-Gemini](https://git.sr.ht/~yotam/go-gemini) and
[Melchior](https://github.com/praetoriansentry/melchior) and, of course, the [Gemini spec](https://gemini.circumlunar.space/docs/specification.html) itself, which is a masterpiece in minimalism.

A week or so ago, I was discovering Mastadon and ActivityPub, and thinking about building an ActivityPub
server. Then, Gemini swept in and I was truly inspired. The focus on security, the minimalism, the philosophy
and the DIY ethic are a joy to behold, not to mention the aesthetic; Gemini capsules in Geminispace, with
the protocol positioned somewhere between Mercury (Gopher) and Apollo (the Web). Great stuff.

In my little corner of cyberspace, there has been so much attention lavished on Gemini so quickly
that I can't now remember (it's only been a week!) exactly where I heard/read about it first, but it would be
remiss of me not to mention Chris Were, HexDSL and Uoou, who have produced some of the best videos on this 
subject.

# Specification, I don't need no stinking specification

Seriously, the spec is obviously vitally important, and none of the fun that people are having with servers
and browsers and producing content would be possible without the work of the authors of the Gemini protocol.

However, minimal as it is, there are parts of the spec that do not apply to my use case. If I am missing
anything important, however, do please let me know.

## Input

I don't envision using it for anything, for the time being, and to be honest I'm not sure I understand
what the spec is trying to say.

## Redirects

Likewise, I don't envision needing a server that implements redirects for my purposes, and again, I'm
not entirely sure how it would work.

## Client Certs

The spec itself acknowledges that minimal implementations are unlikely to make use of this capability.
