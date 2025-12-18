<div align="center">
  <img src="header.png" alt="TempMailChecker Go SDK" width="100%">
</div>

# TempMailChecker Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/Fushey/go-disposable-email-checker.svg)](https://pkg.go.dev/github.com/Fushey/go-disposable-email-checker)
[![Go Report Card](https://goreportcard.com/badge/github.com/Fushey/go-disposable-email-checker)](https://goreportcard.com/report/github.com/Fushey/go-disposable-email-checker)
[![License](https://img.shields.io/github/license/Fushey/go-disposable-email-checker?style=flat-square)](LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/Fushey/go-disposable-email-checker?style=flat-square)](https://github.com/Fushey/go-disposable-email-checker/stargazers)

> **Detect disposable email addresses in real-time** with the TempMailChecker API. Block fake signups, prevent spam, and protect your platform from abuse. Zero dependencies.

## 🚀 Quick Start

### Installation

```bash
go get github.com/Fushey/go-disposable-email-checker
```

### Basic Usage

```go
package main

import (
    "fmt"
    "log"

    tempmailchecker "github.com/Fushey/go-disposable-email-checker"
)

func main() {
    // Create client
    checker, err := tempmailchecker.New("your_api_key")
    if err != nil {
        log.Fatal(err)
    }

    // Check an email
    result, err := checker.Check("user@tempmail.com")
    if err != nil {
        log.Fatal(err)
    }

    if result.Temp {
        fmt.Println("⚠️  Disposable email detected!")
    } else {
        fmt.Println("✅ Legitimate email")
    }
}
```

## 📖 API Reference

### Creating a Client

```go
// Basic client with default settings (EU endpoint)
checker, err := tempmailchecker.New("your_api_key")

// With custom endpoint for lower latency
checker, err := tempmailchecker.New("your_api_key",
    tempmailchecker.WithEndpoint(tempmailchecker.EndpointUS),
)

// With custom timeout
checker, err := tempmailchecker.New("your_api_key",
    tempmailchecker.WithTimeout(5 * time.Second),
)

// Panic on error (for initialization)
checker := tempmailchecker.MustNew("your_api_key")
```

### Regional Endpoints

Choose the endpoint closest to your users for lowest latency:

| Constant | Region | Best For |
|----------|--------|----------|
| `EndpointEU` | 🇪🇺 Europe | EU, Africa, Middle East |
| `EndpointUS` | 🇺🇸 United States | Americas |
| `EndpointAsia` | 🇸🇬 Asia | Asia-Pacific, Australia, Japan |

```go
// US endpoint
checker, _ := tempmailchecker.New("your_api_key",
    tempmailchecker.WithEndpoint(tempmailchecker.EndpointUS),
)

// Asia endpoint
checker, _ := tempmailchecker.New("your_api_key",
    tempmailchecker.WithEndpoint(tempmailchecker.EndpointAsia),
)
```

### Check Email

```go
result, err := checker.Check("user@example.com")
if err != nil {
    log.Fatal(err)
}

if result.Temp {
    fmt.Println("Disposable email!")
}
```

### Check Domain

```go
result, err := checker.CheckDomain("tempmail.com")
if err != nil {
    log.Fatal(err)
}

if result.Temp {
    fmt.Println("Disposable domain!")
}
```

### Convenience Method

```go
isDisposable, err := checker.IsDisposable("user@tempmail.com")
if err != nil {
    log.Fatal(err)
}

if isDisposable {
    fmt.Println("Blocked!")
}
```

### Get Usage Statistics

```go
usage, err := checker.GetUsage()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Used %d of %d requests today\n", usage.UsageToday, usage.Limit)
fmt.Printf("Resets at: %s\n", usage.Reset)
```

## 🛡️ Error Handling

```go
result, err := checker.Check(email)
if err != nil {
    // Check for specific error types
    if tempmailchecker.IsRateLimitError(err) {
        fmt.Println("Rate limit exceeded, try again later")
        return
    }
    
    if tempmailchecker.IsAPIError(err) {
        fmt.Printf("API error: %v\n", err)
        return
    }
    
    // Validation errors
    switch err {
    case tempmailchecker.ErrEmailRequired:
        fmt.Println("Email is required")
    case tempmailchecker.ErrInvalidEmail:
        fmt.Println("Invalid email format")
    default:
        fmt.Printf("Error: %v\n", err)
    }
    return
}
```

## 🔧 Complete Example

```go
package main

import (
    "fmt"
    "log"
    "time"

    tempmailchecker "github.com/Fushey/go-disposable-email-checker"
)

func main() {
    // Create client with options
    checker, err := tempmailchecker.New("your_api_key",
        tempmailchecker.WithEndpoint(tempmailchecker.EndpointUS),
        tempmailchecker.WithTimeout(5*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }

    // Emails to check
    emails := []string{
        "user@gmail.com",
        "test@10minutemail.com",
        "hello@tempmail.org",
    }

    for _, email := range emails {
        result, err := checker.Check(email)
        if err != nil {
            if tempmailchecker.IsRateLimitError(err) {
                fmt.Println("Rate limit reached!")
                break
            }
            fmt.Printf("Error checking %s: %v\n", email, err)
            continue
        }

        status := "✅ Legitimate"
        if result.Temp {
            status = "⚠️  Disposable"
        }
        fmt.Printf("%s: %s\n", email, status)
    }

    // Check usage
    usage, err := checker.GetUsage()
    if err == nil {
        fmt.Printf("\nAPI Usage: %d/%d requests today\n", usage.UsageToday, usage.Limit)
    }
}
```

## 📊 Response Types

### CheckResult

```go
type CheckResult struct {
    Temp bool `json:"temp"` // true if disposable
}
```

### UsageResult

```go
type UsageResult struct {
    UsageToday int    `json:"usage_today"` // Requests made today
    Limit      int    `json:"limit"`       // Daily limit
    Reset      string `json:"reset"`       // Reset time
}
```

## 🧪 Testing

```bash
# Run unit tests
go test -v

# Run with integration tests (requires API key)
TEMPMAILCHECKER_API_KEY=your_key go test -v
```

## 📝 License

MIT License - see [LICENSE](LICENSE) for details.

## 🔗 Links

- [TempMailChecker Website](https://tempmailchecker.com)
- [API Documentation](https://tempmailchecker.com/docs)
- [Get Your Free API Key](https://tempmailchecker.com)

