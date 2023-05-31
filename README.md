# terraform-provider-papertrail

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.10.x
-	[Go](https://golang.org/doc/install) 1.8 (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/last9/terraform-provider-papertrail`

```sh
$ mkdir -p $GOPATH/src/github.com/last9; cd $GOPATH/src/github.com/last9
$ git clone git@github.com:last9/terraform-provider-papertrail
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/last9/terraform-provider-papertrail
$ go get
$ go build
```

For Usage, have a look at docs in `website` directory.

Running Tests
-------------
```sh
$ cd $GOPATH/src/github.com/last9/terraform-provider-papertrail/papertrail
$ PAPERTRAIL_TOKEN=<token> DESTINATION_PORT=<log_destination_port> go tests -v
```
