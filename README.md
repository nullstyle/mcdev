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

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## License

See [LICENSE.txt](LICENSE.txt)
