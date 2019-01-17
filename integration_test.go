package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	assert.NoError(t, err, "Error while POSTing data")
	res, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode, "Unexpected status code")
}

func TestMain(m *testing.M) {
	flag.Parse()
	if *integrationTest {
		m.Run()
	}
}
