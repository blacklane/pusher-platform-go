package pusherplatform

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestInstanceConstruction(t *testing.T) {
	t.Run("Incorrect instance locator format", func(t *testing.T) {
		_, err := NewInstance(InstanceOptions{
			Locator: "invalid-locator",
		})
		if err != nil {
			if err.Error() != "Instance locator must be of the format <version>:<cluster>:<instance-id>" {
				t.Fatalf("Expected incorrect instance locator error, but got %+v", err)
			}
		}
	})

	t.Run("Incorrect key format", func(t *testing.T) {
		_, err := NewInstance(InstanceOptions{
			Locator: "v1:local:instance-id",
			Key:     "blah",
		})
		if err != nil {
			if err.Error() != "Key must be of the format <key>:<secret>" {
				t.Fatalf("Expected incorrect key error, but got %+v", err)
			}
		}
	})
}

func TestInstanceRequestSuccess(t *testing.T) {
	instanceLocator := "v1:local:instance-id"
	jwt := "jwt"

	mux := http.NewServeMux()
	mux.HandleFunc("/services/test_service/v1/instance-id/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	uri, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server URL: %+v", err)
	}

	baseClient := NewBaseClient(BaseClientOptions{
		Host: uri.Host,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	instance, err := NewInstance(InstanceOptions{
		Locator:        instanceLocator,
		Key:            "key:secret",
		ServiceName:    "test_service",
		ServiceVersion: "v1",
		Client:         baseClient,
	})
	if err != nil {
		t.Fatalf("Expected no error when constructing an instance, but got %+v", err)
	}

	response, err := instance.Request(context.Background(), RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
		Jwt:    &jwt,
	})
	if err != nil {
		t.Fatalf("Expected no error when performing a request, but got %+v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected a 200 status code, but got %v", response.StatusCode)
	}
}

func TestInstanceRequestSuccessWithoutJwt(t *testing.T) {
	instanceLocator := "v1:local:instance-id"

	mux := http.NewServeMux()
	mux.HandleFunc("/services/test_service/v1/instance-id/test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := httptest.NewTLSServer(mux)
	defer server.Close()

	uri, err := url.Parse(server.URL)
	if err != nil {
		t.Fatalf("Failed to parse server URL: %+v", err)
	}

	baseClient := NewBaseClient(BaseClientOptions{
		Host: uri.Host,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	})

	instance, err := NewInstance(InstanceOptions{
		Locator:        instanceLocator,
		Key:            "key:secret",
		ServiceName:    "test_service",
		ServiceVersion: "v1",
		Client:         baseClient,
	})
	if err != nil {
		t.Fatalf("Expected no error when constructing an instance, but got %+v", err)
	}

	response, err := instance.Request(context.Background(), RequestOptions{
		Method: http.MethodGet,
		Path:   "/test",
	})
	if err != nil {
		t.Fatalf("Expected no error when performing a request, but got %+v", err)
	}

	if response.StatusCode != http.StatusOK {
		t.Fatalf("Expected a 200 status code, but got %v", response.StatusCode)
	}
}
