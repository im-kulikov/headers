# Unmarshal HTTP Headers into Struct Fields

[![Go Reference](https://pkg.go.dev/badge/github.com/im-kulikov/headers.svg)](https://pkg.go.dev/github.com/im-kulikov/headers)
[![Go Report Card](https://goreportcard.com/badge/github.com/im-kulikov/headers)](https://goreportcard.com/report/github.com/im-kulikov/headers)

## Overview

The `headers` Go package provides an easy and efficient way to unmarshal HTTP headers into Go struct fields using struct tags and reflection. It supports a wide range of data types and custom deserialization logic via the `CustomUnmarshaler` interface.

### Features

- **Automatic Header Mapping**: Unmarshals HTTP headers into struct fields using the `header` struct tag.
- **Support for Standard Go Types**: Handles various types like `string`, `bool`, `int`, `uint`, `float`, and slices.
- **Custom Unmarshaling**: Implement the `CustomUnmarshaler` interface for custom parsing logic.
- **Error Handling**: Returns clear errors for unknown types or invalid header values.

### Installation

Install the package using `go get`:

```bash
go get github.com/im-kulikov/headers
```

### Usage

Define your struct with fields mapped to HTTP headers using the `header` struct tag:

```go
package main

import (
    "fmt"
    "net/http"
	
    "github.com/im-kulikov/headers"
)

type MyHeaders struct {
    AuthToken string `header:"Authorization"`
    UserID    int    `header:"X-User-ID"`
    IsAdmin   bool   `header:"X-Is-Admin"`
}

func main() {
    hdr := http.Header{}
    hdr.Add("Authorization", "Bearer some-token")
    hdr.Add("X-User-ID", "123")
    hdr.Add("X-Is-Admin", "true")

    var myHeaders MyHeaders
    err := headers.UnmarshalHeaders(&myHeaders, hdr)
    if err != nil {
        fmt.Println("Error:", err)
    }

    fmt.Printf("Parsed Headers: %+v\n", myHeaders)
}
```

### Custom Unmarshaling

You can define custom unmarshaling logic by implementing the `CustomUnmarshaler` interface. This is useful when you need to handle complex types or perform additional validation.

```go
type CustomType struct {
    Value string
}

// Custom unmarshaling logic
func (c *CustomType) UnmarshalHeader(values []string) error {
    if len(values) > 0 {
        c.Value = "custom:" + values[0]
    }
    return nil
}

type MyHeaders struct {
    CustomField CustomType `header:"X-Custom-Field"`
}
```

### Struct Tag Behavior

- **`header:"<Header-Name>"`**: Binds the header to the struct field.
- **`header:"-"`**: Skips unmarshaling for the field.

### Error Handling

If an unknown type is encountered or parsing fails (e.g., invalid integer or boolean), the function returns an error, making it easy to handle invalid headers in your application.

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.