package tempmailchecker

import (
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("valid API key", func(t *testing.T) {
		client, err := New("test_key")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if client == nil {
			t.Fatal("expected client, got nil")
		}
	})

	t.Run("empty API key", func(t *testing.T) {
		_, err := New("")
		if err != ErrAPIKeyRequired {
			t.Fatalf("expected ErrAPIKeyRequired, got %v", err)
		}
	})

	t.Run("whitespace API key", func(t *testing.T) {
		_, err := New("   ")
		if err != ErrAPIKeyRequired {
			t.Fatalf("expected ErrAPIKeyRequired, got %v", err)
		}
	})
}

func TestMustNew(t *testing.T) {
	t.Run("valid API key", func(t *testing.T) {
		client := MustNew("test_key")
		if client == nil {
			t.Fatal("expected client, got nil")
		}
	})

	t.Run("empty API key panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Fatal("expected panic, got none")
			}
		}()
		MustNew("")
	})
}

func TestWithOptions(t *testing.T) {
	t.Run("with endpoint", func(t *testing.T) {
		client, _ := New("test_key", WithEndpoint(EndpointUS))
		if client.endpoint != EndpointUS {
			t.Fatalf("expected %s, got %s", EndpointUS, client.endpoint)
		}
	})

	t.Run("with timeout", func(t *testing.T) {
		client, _ := New("test_key", WithTimeout(5*time.Second))
		if client.timeout != 5*time.Second {
			t.Fatalf("expected 5s, got %v", client.timeout)
		}
	})
}

func TestEndpointConstants(t *testing.T) {
	if EndpointEU != "https://tempmailchecker.com" {
		t.Errorf("EndpointEU = %s; want https://tempmailchecker.com", EndpointEU)
	}
	if EndpointUS != "https://us.tempmailchecker.com" {
		t.Errorf("EndpointUS = %s; want https://us.tempmailchecker.com", EndpointUS)
	}
	if EndpointAsia != "https://asia.tempmailchecker.com" {
		t.Errorf("EndpointAsia = %s; want https://asia.tempmailchecker.com", EndpointAsia)
	}
}

func TestCheckValidation(t *testing.T) {
	client := MustNew("test_key")

	t.Run("empty email", func(t *testing.T) {
		_, err := client.Check("")
		if err != ErrEmailRequired {
			t.Fatalf("expected ErrEmailRequired, got %v", err)
		}
	})

	t.Run("whitespace email", func(t *testing.T) {
		_, err := client.Check("   ")
		if err != ErrEmailRequired {
			t.Fatalf("expected ErrEmailRequired, got %v", err)
		}
	})

	t.Run("invalid email format", func(t *testing.T) {
		_, err := client.Check("not-an-email")
		if err != ErrInvalidEmail {
			t.Fatalf("expected ErrInvalidEmail, got %v", err)
		}
	})
}

func TestCheckDomainValidation(t *testing.T) {
	client := MustNew("test_key")

	t.Run("empty domain", func(t *testing.T) {
		_, err := client.CheckDomain("")
		if err != ErrDomainRequired {
			t.Fatalf("expected ErrDomainRequired, got %v", err)
		}
	})

	t.Run("whitespace domain", func(t *testing.T) {
		_, err := client.CheckDomain("   ")
		if err != ErrDomainRequired {
			t.Fatalf("expected ErrDomainRequired, got %v", err)
		}
	})
}

func TestErrorTypes(t *testing.T) {
	t.Run("IsRateLimitError", func(t *testing.T) {
		err := &RateLimitError{Message: "test"}
		if !IsRateLimitError(err) {
			t.Fatal("expected IsRateLimitError to return true")
		}
		if IsRateLimitError(ErrAPIKeyRequired) {
			t.Fatal("expected IsRateLimitError to return false for non-rate-limit error")
		}
	})

	t.Run("IsAPIError", func(t *testing.T) {
		err := &APIError{StatusCode: 400, Message: "test"}
		if !IsAPIError(err) {
			t.Fatal("expected IsAPIError to return true")
		}
		if IsAPIError(ErrAPIKeyRequired) {
			t.Fatal("expected IsAPIError to return false for non-API error")
		}
	})
}

// Integration tests - only run when API key is available
func TestIntegration(t *testing.T) {
	apiKey := os.Getenv("TEMPMAILCHECKER_API_KEY")
	if apiKey == "" {
		t.Skip("TEMPMAILCHECKER_API_KEY not set, skipping integration tests")
	}

	client := MustNew(apiKey)

	t.Run("check disposable email", func(t *testing.T) {
		result, err := client.Check("test@10minutemail.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Temp {
			t.Error("expected Temp=true for disposable email")
		}
	})

	t.Run("check legitimate email", func(t *testing.T) {
		result, err := client.Check("test@gmail.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result.Temp {
			t.Error("expected Temp=false for legitimate email")
		}
	})

	t.Run("is disposable", func(t *testing.T) {
		isDisposable, err := client.IsDisposable("test@tempmail.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !isDisposable {
			t.Error("expected IsDisposable=true")
		}
	})

	t.Run("check domain", func(t *testing.T) {
		result, err := client.CheckDomain("tempmail.com")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !result.Temp {
			t.Error("expected Temp=true for disposable domain")
		}
	})

	t.Run("get usage", func(t *testing.T) {
		usage, err := client.GetUsage()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if usage.Limit <= 0 {
			t.Error("expected Limit > 0")
		}
	})
}

