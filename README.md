# Mcdev, a suite of tools to spice up your go development workflow.


## Installation

Installation requires go 1.5 with the go 1.5 vendoring experiment turned on by
setting the environment variable `GO15VENDOREXPERIMENT=1`.  Given that:

```
go get github.com/nullstyle/mcdev/cmd/...
```

## Usage

### Run a package's test everytime one of its .go files change
```
mcdev-each-change go test {{.Pkg}}
```

### Stop and re-start the server any time a package underneath the pwd is changed
```
mcdev-rerun go run examples/server.go
```

### gb mode

The `pkgwatch` package converts notifications of changed files into
notifications of changed go packages.  Normally, this is done by searching the
GOPATH for parents of the file that was changed.  This doesn't play well with
the [gb tool](http://getgb.io/).

Instead, we provide a "gb mode" for the pkgwatch package, enabled by setting the
`-gb` flag on the command line.  When enabled and the tools are run from a gb
project root, the correct package changes will be detected.

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## License

See [LICENSE.txt](LICENSE.txt)
