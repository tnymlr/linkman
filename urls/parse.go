package urls

import (
	"fmt"
	"net/url"

	tld "golang.org/x/net/publicsuffix"
)

func GetSource(url *url.URL) (string, error) {
	msg := "Unable to create source string: %"

	if err := validate(url); err != nil {
		return "", fmt.Errorf(msg, err)
	} else {
		return extractSource(url.Hostname(), msg)
	}
}

func ParseUrl(rawurl string) (*url.URL, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	} else {
		return url, validate(url)
	}
}

func extractSource(hostname string, errMsg string) (string, error) {
	if tldPlusOne, err := tld.EffectiveTLDPlusOne(hostname); err == nil {
		suffix, _ := tld.PublicSuffix(tldPlusOne)
		return tldPlusOne[:len(tldPlusOne)-len(suffix)-1], nil
	} else {
		return "", fmt.Errorf(errMsg, err)
	}

}

func validate(subjectUrl *url.URL) error {
	var err error
	validators := []func(*url.URL) error{
		validateHost,
	}

	for _, v := range validators {
		err = v(subjectUrl)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateHost(url *url.URL) error {
	if url.Host == "" || url.Hostname() == "" {
		return fmt.Errorf("Unexpected URL: missing hostname")
	} else {
		return nil
	}
}
