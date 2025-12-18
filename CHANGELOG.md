# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-12-18

### Added
- Initial release
- `New()` and `MustNew()` client constructors
- `Check()` method to verify email addresses
- `CheckDomain()` method to verify domains
- `IsDisposable()` convenience method
- `GetUsage()` method to check API usage
- Regional endpoint support (EU, US, Asia)
- Configurable timeout via `WithTimeout()`
- Custom HTTP client support via `WithHTTPClient()`
- Comprehensive error types (`APIError`, `RateLimitError`)
- Helper functions `IsRateLimitError()` and `IsAPIError()`
- Full test coverage
- Examples and documentation

