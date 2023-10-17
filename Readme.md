# Slot

[![Go Reference](https://pkg.go.dev/badge/github.com/livebud/slot.svg)](https://pkg.go.dev/github.com/livebud/slot)

Compose a response from nested HTTP handlers. A generic version of [Svelte Slots](https://svelte.dev/examples/slots). Used by [mux](http://github.com/livebud/mux).

Supports handlers running in series or concurrently (like [Remix](https://remix.run/docs/en/main/discussion/routes)).

## Install

```sh
go get -u github.com/livebud/slot
```

## Example

See `ExampleBatch` in [slot_test.go](./slot_test.go).

## Contributors

- Matt Mueller ([@mattmueller](https://twitter.com/mattmueller))

## License

MIT
