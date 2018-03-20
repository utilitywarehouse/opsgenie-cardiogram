package cardiogram

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
		APIUrl := fmt.Sprintf("https://api.opsgenie.com/v2/heartbeats/%s/ping", name)
		h.send(APIUrl)
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

func (h *Heartbeat) send(APIUrl string) {
	req, err := http.NewRequest("POST", APIUrl, nil)
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return
	}

	apiKey := fmt.Sprintf("GenieKey %s", h.APIKey)
	req.Header.Set("Authorization", apiKey)

	res, err := h.Client.Do(req)
	if err != nil {
		log.Printf("Error sending the Heartbeat request: %s", err)
		return
	}

	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	if res.StatusCode != 202 {
		reply, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Printf("Error reading opsgenie reply body: %s", err)
			return
		}
		log.Println("Opsgenie reply to Heartbeat not successful: " + string(reply))
	}
}
