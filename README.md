# Clacks

A tiny Go HTTP middleware for adding the X-Clacks-Overhead header to HTTP responses, keeping the memory of Terry Pratchett alive on the internet.

## About

> "A man is not dead while his name is still spoken."
> — Terry Pratchett, Going Postal

This middleware implements the [GNU Terry Pratchett](https://xclacksoverhead.org/) protocol, adding the `X-Clacks-Overhead: GNU Terry Pratchett` header to all HTTP responses.

In Terry Pratchett's Discworld series, the Clacks is a semaphore system for sending messages across long distances. When Robert Dearheart's son John died, his name was inserted into the overhead of the Clacks with a "GNU" code, causing his name to be repeated forever throughout the network.

## Installation

```bash
go get github.com/flaticols/clacks
```

## Usage

### With http.ServeMux

```go
package main

import (
    "net/http"
    "github.com/flaticols/clacks"
)

func main() {
    mux := http.NewServeMux()

    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    // Wrap your mux with the Clacks middleware
    handler := clacks.Clacks(mux)

    http.ListenAndServe(":8080", handler)
}
```

### With any http.Handler

```go
package main

import (
    "net/http"
    "github.com/flaticols/clacks"
)

func main() {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    // Wrap your handler
    wrappedHandler := clacks.Middleware(handler)

    http.ListenAndServe(":8080", wrappedHandler)
}
```

### With go-pkgz/routegroup

```go
package main

import (
    "net/http"
    "github.com/flaticols/clacks"
    "github.com/go-pkgz/routegroup"
)

func main() {
    // Create a new route group with Clacks middleware
    rg := routegroup.New(http.NewServeMux())
    rg.Use(clacks.Clacks)

    // Add routes to the group
    rg.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })

    rg.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("About page"))
    })

    http.ListenAndServe(":8080", rg)
}
```

You can also apply the middleware to specific route groups:

```go
package main

import (
    "net/http"
    "github.com/flaticols/clacks"
    "github.com/go-pkgz/routegroup"
)

func main() {
    rg := routegroup.New(http.NewServeMux())

    // API routes with Clacks middleware
    api := rg.Mount("/api")
    api.Use(clacks.Clacks)
    api.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Users API"))
    })

    // Public routes without middleware
    rg.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Home page"))
    })

    http.ListenAndServe(":8080", rg)
}
```

## Testing

```bash
go test -v
```

## How It Works

The middleware wraps your HTTP handler and adds the `X-Clacks-Overhead: GNU Terry Pratchett` header to every response before passing control to your handler.

The GNU overhead code means:

- **G**: Send the message on
- **N**: Do not log the message
- **U**: Turn the message around at the end of the line

## License

MIT

## References

- [X-Clacks-Overhead](https://xclacksoverhead.org/)
- [GNU Terry Pratchett](http://www.gnuterrypratchett.com/)
