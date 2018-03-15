package cardiogram

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Heartbeat contains the configuration for Opsgenie Heartbeats.
type Heartbeat struct {
	Client  *http.Client
	Timeout time.Duration
	APIKey  string
}

// Check scrapes the targets and send the heartbeats to Opsgenie.
func (h *Heartbeat) Check(url string, expected int, name string) {
	if h.call(url, expected) == nil {
		h.send(name)
	}
}

func (h *Heartbeat) call(url string, expected int) error {
	res, err := h.Client.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != expected {
		return errors.New("Target returns an unexpected status code")
	}
	return nil
}

func (h *Heartbeat) send(name string) {
	url := fmt.Sprintf("https://api.opsgenie.com/v2/heartbeats/%s/ping", name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Error creating request: %s", err)
	}

	apiKey := fmt.Sprintf("GenieKey %s", h.APIKey)
	req.Header.Set("Authorization", apiKey)

	res, err := h.Client.Do(req)
	if err != nil {
		log.Printf("Error while sending Heartbeat for '%s': %s", name, err)
	}

	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	if res.StatusCode != 202 {
		log.Println("Sending Heartbeat was not successful")
	}
}
