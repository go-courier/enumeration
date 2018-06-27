## Enumeration

[![GoDoc Widget](https://godoc.org/github.com/go-courier/enumeration?status.svg)](https://godoc.org/github.com/go-courier/enumeration)
[![Build Status](https://travis-ci.org/go-courier/enumeration.svg?branch=master)](https://travis-ci.org/go-courier/enumeration)
[![codecov](https://codecov.io/gh/go-courier/enumeration/branch/master/graph/badge.svg)](https://codecov.io/gh/go-courier/enumeration)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-courier/enumeration)](https://goreportcard.com/report/github.com/go-courier/enumeration)

Enumeration in go


```go
// declare named type base on int or uint
type Protocol int

// declare const as the declared type
// and with prefix upper-snake-cased type name
// link with UNKNOWN by one lodash _ as the zero value
// link with ENUM VALUE by two lodash __ as enum values
// comments after enum values will be the label of matched enum value
const (
	PROTOCOL_UNKNOWN Protocol = iota
	PROTOCOL__HTTP    // http
	PROTOCOL__HTTPS   // https
	PROTOCOL__TCP
)

```