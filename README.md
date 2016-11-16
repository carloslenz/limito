limito
======

Go library to limit concurrent operations.

```sh
go get github.com/carloslenz/limito
```

Available abstractions are:

* `Limiter`: limits concurrent operations.
* `WaitList`: serializes access to resource.
* `Middleware`: restricts concurrent HTTP requests per user ID.
**Important**: Set a valid ID (`SetMiddlewareID`) in an upper layer.
If you don't, `Middleware` will `panic` to let you know there is a bug!

License
-------

MIT