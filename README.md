# ecosystem

Ecosystem includes tooling that helps interact with xo/xo in various ways.

## protoc-gen-xo

protoc-gen-xo is a xo types generator, generating an sql schema from protobuf
definitions. It is designed to be used with gunk, but can work separately,
outside of gunk.

### Installation

Use the following command to install protoc-gen-xo:

```sh
$ go install github.com/xo/ecosystem/cmd/protoc-gen-xo@latest
```
