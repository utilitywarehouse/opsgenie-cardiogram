package cardiogram

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockServer(statusCode int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, body)
	}))
}

func TestSend(t *testing.T) {
	h := Heartbeat{
		Client: &http.Client{},
		APIKey: "1234",
	}
	t.Run("Invalid URL", func(t *testing.T) {
		err := h.send("http://example.com/%%%%")
		if !strings.Contains(err.Error(), "Error creating request") {
			t.Error("Invalid URL not detected")
		}
	})
	t.Run("Hearbeat not successful", func(t *testing.T) {
		errMsg := "Reply error message"
		server := mockServer(400, errMsg)
		defer server.Close()

		err := h.send(server.URL)
		if !strings.Contains(err.Error(), "Opsgenie reply to Heartbeat not successful") ||
			!strings.Contains(err.Error(), errMsg) {
			t.Fail()
		}
	})
	t.Run("Hearbeat successful", func(t *testing.T) {
		server := mockServer(202, "")
		defer server.Close()
		err := h.send(server.URL)
		if err != nil {
			t.Fail()
		}
	})
}
