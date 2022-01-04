# curly

[![Go Report Card](https://goreportcard.com/badge/github.com/bilal-bhatti/curly)](https://goreportcard.com/report/github.com/bilal-bhatti/curly)


JSON API request collections.

## install
``` sh
brew tap bilal-bhatti/homebrew-taps
brew install curly
```
or
``` sh
go install github.com/bilal-bhatti/curly/cmd/curly@latest
```

## commands
get started

``` sh
curly --help
```

### run
make an api call

``` sh
curly --verbose httpbin/get.status.yml
curly httpbin/post.something.yml
eval "$(curly -c httpbin/post.something.yml)"
```

