## Writing filters

Filters must conform to the interface:

```go
type Filter interface {
	Execute(ctx context.Context, inputTransaction fetch.Transaction) (next bool, err error)
}
```

See examples in the `internal/filter` folder.
