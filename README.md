# pagecache

[![Go Documentation](https://godocs.io/git.sr.ht/~jamesponddotco/pagecache-go?status.svg)](https://godocs.io/git.sr.ht/~jamesponddotco/pagecache-go)
[![Go Report Card](https://goreportcard.com/badge/git.sr.ht/~jamesponddotco/pagecache-go)](https://goreportcard.com/report/git.sr.ht/~jamesponddotco/pagecache-go)
[![pagecache-go build status](https://builds.sr.ht/~jamesponddotco/pagecache-go.svg)](https://builds.sr.ht/~jamesponddotco/pagecache-go?)

Package `pagecache` provides a stable interface and helpers for caching HTTP responses.

## Features

- Stable cache interface.
- Simple and easy-to-use API.
- Multiple helpers, making implementation easier.

### `pagecache.Cache` implementations

The `pagecace` package itself only provides a cache interface and helper
functions for users who wish to implement that interface. You can either
use an implementation created by someone else or write your own.

**Implementations**

- [`memorycachex`](https://git.sr.ht/~jamesponddotco/pagecache-go/tree/trunk/item/memorycachex)
  provides a thread-safe in-memory cache using the
  [Mockingjay](https://en.wikipedia.org/wiki/Cache_replacement_policies#Mockingjay)
  cache replacement policy.

If you wrote a `pagecache.Cache` implementation and wish it to be linked
here, [please send a patch](https://git.sr.ht/~jamesponddotco/pagecache-go#resources).

## Installation

To install `pagecache` alone, run:

```sh
go get git.sr.ht/~jamesponddotco/pagecache-go
```

## Contributing

Anyone can help make `pagecache` better. Check out [the contribution
guidelines](https://git.sr.ht/~jamesponddotco/pagecache-go/tree/master/item/CONTRIBUTING.md)
for more information.

## Resources

The following resources are available:

- [Package documentation](https://godocs.io/git.sr.ht/~jamesponddotco/pagecache-go).
- [Support and general discussions](https://lists.sr.ht/~jamesponddotco/pagecache-discuss).
- [Patches and development related questions](https://lists.sr.ht/~jamesponddotco/pagecache-devel).
- [Instructions on how to prepare patches](https://git-send-email.io/).
- [Feature requests and bug reports](https://todo.sr.ht/~jamesponddotco/pagecache).

---

Released under the [MIT License](LICENSE.md).
