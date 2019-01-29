package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	try "gopkg.in/matryer/try.v1"
)

func resolveIP(provider string) (string, error) {
	const maxAttempts = 5
	var resolvedIP string
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		resolvedIP, err = tryResolveIP(provider)
		return attempt < maxAttempts, err
	})

	return resolvedIP, err
}

func tryResolveIP(provider string) (string, error) {
	resp, err := http.Get(provider)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Received non 200 status code '%v' from '%v'", resp.StatusCode, provider)
	}
	resolvedIP := strings.TrimSpace(string(bodyBytes))
	if net.ParseIP(resolvedIP) == nil {
		return "", fmt.Errorf("Response invalid, unable to parse IP")
	}
	return resolvedIP, nil
}
