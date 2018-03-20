package cardiogram

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
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
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		h.send("http://example.com/%%%%")
		log := buf.String()
		if !strings.Contains(log, "Error creating request") {
			t.Error("Invalid URL not detected")
		}
	})
	t.Run("Hearbeat not successful", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		errMsg := "Reply error message"
		server := mockServer(400, errMsg)
		defer server.Close()

		h.send(server.URL)
		log := buf.String()
		if !strings.Contains(log, "Opsgenie reply to Heartbeat not successful") {
			t.Fail()
		}
		if !strings.Contains(log, errMsg) {
			t.Fail()
		}
	})
	t.Run("Hearbeat successful", func(t *testing.T) {
		var buf bytes.Buffer
		log.SetOutput(&buf)
		defer func() {
			log.SetOutput(os.Stderr)
		}()

		server := mockServer(202, "")
		defer server.Close()

		h.send(server.URL)
		log := buf.String()
		if strings.Contains(log, "Opsgenie reply to Heartbeat not successful") {
			t.Fail()
		}
	})
}
