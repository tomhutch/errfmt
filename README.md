# Errfmt

A Golang linter that checks whether wrapped errors have a consistent format.

## Installation

```sh
go install github.com/tomhutch/errfmt
```

## Usage

The `errfmt` linter is called using the following format: `errfmt [-flag] [package]`.
Where `[package]` can be a Golang package, a filepath or `./...` for all files recursively, e.g.:
```sh
errfmt ./...
```

Apply suggested fixes using the `-fix` flag:
```sh
errfmt -fix ./...
```
