# Contributing to TempMailChecker Go SDK

Thank you for your interest in contributing! This document provides guidelines and instructions for contributing to this project.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/go-disposable-email-checker.git`
3. Create a new branch: `git checkout -b feature/your-feature-name`

## Development

### Prerequisites

- Go 1.18 or later
- Make (optional, for convenience commands)

### Running Tests

```bash
# Run unit tests
go test -v

# Run with coverage
go test -v -cover

# Run integration tests (requires API key)
TEMPMAILCHECKER_API_KEY=your_key go test -v
```

### Code Style

- Follow standard Go conventions and idioms
- Run `go fmt` before committing
- Run `go vet` to catch common issues
- Use meaningful variable and function names
- Add comments for exported functions and types

### Building

```bash
go build ./...
```

## Submitting Changes

1. Ensure all tests pass: `go test -v`
2. Run `go fmt` and `go vet`
3. Commit your changes with clear, descriptive messages
4. Push to your fork: `git push origin feature/your-feature-name`
5. Create a Pull Request on GitHub
6. Provide a clear description of your changes

## Reporting Issues

When reporting issues, please include:
- Go version (`go version`)
- Operating system
- Steps to reproduce
- Expected vs actual behavior
- Any error messages

## Questions?

Feel free to open an issue for questions or reach out to support@tempmailchecker.com

Thank you for contributing! 🎉

