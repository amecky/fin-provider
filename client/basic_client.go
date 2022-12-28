package client

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func LoadHtmlContent(url string) (string, error) {
	httpClient := http.Client{
		Timeout: time.Second * 4, // Timeout after 2 seconds
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Stocks App")

	res, getErr := httpClient.Do(req)
	if getErr != nil {
		return "", getErr
	}

	if res.StatusCode != 200 {
		fmt.Println("Wrong status code", res.StatusCode)
		return "", errors.New(fmt.Sprintf("Status Code: %d", res.StatusCode))
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return "", readErr
	}
	if err != nil {
		return "", err
	}
	return string(body), nil
}
