My attempt at building a basic version of [`grep`](https://en.wikipedia.org/wiki/Grep), a CLI tool for searching using Regexes.

I built this to learn about Regex, how they are evaluated and how parsers/lexers work.

How to run:

1. Clone the repository
2. Ensure you have `go (1.22)` installed locally
3. Run `echo <text> | go run cmd/mygrep/main.go -E "<pattern>"` to run your program, which is implemented in
   `cmd/mygrep/main.go`.
