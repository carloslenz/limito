limito
======

Go library to limit concurrent operations.

```sh
go get github.com/carloslenz/limito
```

Available abstractions are:

* `Limiter`: limits number of concurrent operations.
* `WaitList`: serializes access to resource (e.g, external service).
* `Middleware`: restricts concurrent HTTP requests per user ID.
**Important**: Set a valid ID (`SetMiddlewareID`) into the context passed
to the `Handler` returned by `Middleware`, otherwise it will `panic` to let
you know there is a bug!

License
-------

MIT