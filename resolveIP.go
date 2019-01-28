package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	try "gopkg.in/matryer/try.v1"
)

func resolveIP(provider string) (string, error) {
	const maxAttempts = 5
	var resolvedIP string
	err := try.Do(func(attempt int) (bool, error) {
		var err error
		resp, err := http.Get(provider)
		if err == nil {
			defer resp.Body.Close()
			bodyBytes, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				if resp.StatusCode != http.StatusOK {
					err = fmt.Errorf("Received non 200 status code '%v' from '%v'", resp.StatusCode, provider)
				} else {
					resolvedIP = strings.TrimSpace(string(bodyBytes))
				}
			}
		}
		if err != nil {
			log.Print(fmt.Sprintf("WARNING: Resolve attempt %v/%v: %v", attempt, maxAttempts, err))
			if attempt != maxAttempts {
				time.Sleep(time.Duration(attempt*attempt) * time.Second)
			}
		}
		return attempt < maxAttempts, err
	})
	return resolvedIP, err
}
