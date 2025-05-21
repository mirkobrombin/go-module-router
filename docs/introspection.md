# Route Introspection

Because the routing table is constructed **before** the server starts, you can
inspect or export it for:

*   Unit tests – assert that critical routes exist.
*   API documentation generators.
*   Health checkers.
*   Static analysis tooling.
*   ... and so on.

```go
rt := registry.Global()

for _, rp := range rt.RouteProviders {
    for _, r := range rp() {
        fmt.Printf("%s %s → %s\n", r.Method, r.Path, r.HandlerName)
    }
}
```

The snippet works without starting an `Engine`, so it is
safe to run in `go test`.
