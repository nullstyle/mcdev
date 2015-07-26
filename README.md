# Mcdev, a suite of tools to spice up your go development workflow.



## Installation

Installation requires go 1.5 with the go 1.5 vendoring experiment turned on by
setting the environment variable `GO15VENDOREXPERIMENT=1`.  Given that:

```
go get github.com/nullstyle/mcdev/cmd/...
```

## Usage

### Run a package's test everytime one if its .go files change
```
mcdev-each-change -cmd='go test {{.Pkg}}'
```

### Stop and re-start the server any time a package underneath the pwd is changed
```
mcdev-rerun examples/server.go
```

## Contributing

## License
