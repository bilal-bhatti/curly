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
curly commands
curly help run
```

### run
make an api call

``` sh
curly -d run httpbin/get.status.yml
curly httpbin/post.something.yml
```

