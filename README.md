api
===

[![Build Status](https://ci.iptq.io/api/badges/nso/api/status.svg)](https://ci.iptq.io/nso/api)

This is the API server, written using the Echo framework for Go.

Building the Project
--------------------

Step 1: Make sure this is cloned to the appropriate location in your `$GOPATH`. Since this project has many moving parts, I decided it would be best to rely on `$GOPATH` instead of using the more recent Go 1.11 modules.

Step 2: Build it as you would a regular Go project.

```
$ go get -v ./...
$ go build ./cmd/api    # this builds the binary
```

Usage
-----

Currently, running this project depends on the existence of `api.yml` in the current working directory. As the project matures, more configuration options will be added. See `config.go` for the approximate structure of the yml file.

Contact
-------

Author: Michael Zhang

License: MIT
