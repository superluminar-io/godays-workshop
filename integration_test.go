package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	endpoint        = flag.String("endpoint", "", "The API gateway endpoint")
	integrationTest = flag.Bool("integrationTest", false, "Run this test suite")
)

// uniqueUrl returns a unique URL
func uniqueUrl() string {
	now := time.Now()
	return fmt.Sprintf("https://godays.io?ts=%s", now)
}

func TestCreateUrl(t *testing.T) {
	// Generate test data
	data, _ := json.Marshal(map[string]string{"url": uniqueUrl()})
	// Perform a HTTP POST request
	client := &http.Client{}
	url := fmt.Sprintf("%s/create-url", *endpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	require.NoError(t, err, "Error while POSTing data")
	res, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode, "Unexpected status code")
}

func TestCreateAndGetUrl(t *testing.T) {
	// Generate test data
	urlToShorten := uniqueUrl()
	data, _ := json.Marshal(map[string]string{"url": urlToShorten})
	// Create the shortened URL
	client := &http.Client{}
	url := fmt.Sprintf("%s/create-url", *endpoint)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	require.NoError(t, err, "Error while POSTing data")
	res, err := client.Do(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode, "Unexpected status code")
	// Retrieve the shortened URL
	body, err := ioutil.ReadAll(res.Body)
	require.NoError(t, err)
	assert.Regexp(t, regexp.MustCompile("^https://"), string(body), "No URL in body")
	shortUrl := string(body)
	res, err = http.Get(shortUrl)
	require.NoError(t, err)
	assert.Equal(t, http.StatusFound, res.StatusCode, "Unexpected status code")
	assert.Equal(t, urlToShorten, res.Location, "Unexpected location")
}

func TestMain(m *testing.M) {
	flag.Parse()
	if *integrationTest {
		m.Run()
	}
}
