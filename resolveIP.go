package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	try "gopkg.in/matryer/try.v1"
)

func resolveIP(provider string) (string, error) {
	const maxAttempts = 5
	var resolvedIP string
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		resolvedIP, err = tryResolveIP(provider)
		if err != nil {
			log.Warn(fmt.Sprintf("Resolve IP attempt %v/%v: %v", attempt, maxAttempts, err))
			if attempt != maxAttempts {
				time.Sleep(time.Duration(attempt*attempt) * time.Second)
			}
		}
		return attempt < maxAttempts, err
	})
	if err != nil {
		log.Error(fmt.Sprintf("Max attempts reached while resolving IP from %v", provider))
	}
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
