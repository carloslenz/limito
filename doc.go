/*
Package limito provides abstractions to limit concurrent operations.

Limiter

Low-level abstraction to limit number of concurrent operations.
New accesses receive errors when the limit is fully in-use.

WaitList

Abstraction that serializes access to resource (e.g, external service).

Middleware

High-level abstraction on top of Limiter. It restricts concurrent HTTP requests per user ID.

IMPORTANT: It requires a valid ID in the request context (otherwise it will panic),
so use SetMiddlewareID in its parent middleware.
*/
package limito
