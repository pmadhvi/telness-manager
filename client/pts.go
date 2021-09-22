package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/pmadhvi/telness-manager/model"
	log "github.com/sirupsen/logrus"
)

type Client struct {
	host string
	log  *log.Logger
}

func NewClient(log *log.Logger, host string) *Client {
	return &Client{
		log:  log,
		host: host,
	}
}

func (c *Client) GetOperatorDetails(msisdn string) (model.PtsResponse, error) {
	// format msisdn number in this format: 010-7500500
	formattedMsisdn := formatMsisdn(msisdn)

	// construct http request to send to pts
	req, err := http.NewRequest(http.MethodGet, c.host, nil)
	if err != nil {
		msg := fmt.Sprintf("could not create http request: %v", err)
		c.log.Errorf(msg)
		return model.PtsResponse{}, err
	}
	req.Header.Add("Accept", "application/json")
	q := req.URL.Query()
	q.Add("Number", string(formattedMsisdn))
	req.URL.RawQuery = q.Encode()

	//define http client with timeout
	client := http.Client{
		Timeout: (1 * time.Second),
	}

	var ptsResponse model.PtsResponse
	// make request to PTS
	response, err := client.Do(req)
	if err != nil {
		msg := fmt.Sprintf("could not get response from PTS: %v", err)
		c.log.Errorf(msg)
		return model.PtsResponse{}, err
	}
	// decode response from PTS to telness model
	err = json.NewDecoder(response.Body).Decode(&ptsResponse)
	if err != nil {
		msg := fmt.Sprintf("could not decode pts response into PtsResponse: %v", err)
		c.log.Error(msg)
		return model.PtsResponse{}, err
	}
	return ptsResponse, nil
}

func formatMsisdn(msisdn string) []rune {
	if strings.Contains(msisdn, "+46") {
		msisdn = strings.Replace(msisdn, "+46", "0", 1)
	}
	var formattedMsisdn []rune
	for i, c := range msisdn {
		if i == 3 {
			formattedMsisdn = append(formattedMsisdn, '-')
		}
		formattedMsisdn = append(formattedMsisdn, c)
	}
	return formattedMsisdn
}
