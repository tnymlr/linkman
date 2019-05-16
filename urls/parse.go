package urls

import (
	"fmt"
	"net/url"

	tld "golang.org/x/net/publicsuffix"
)

//GetSource calculates Source value for provided url.
//Source value is second (or third level) domain name.
//For example:
//| link              | source        |
//| ----              | ------        |
//| youtube.com       | youtube       |
//| stackoverflow.com | stackoverflow |
//| domain.co.uk      | domain        |
func GetSource(url *url.URL) (string, error) {
	msg := "Unable to create source string: %"

	if err := validate(url); err != nil {
		return "", fmt.Errorf(msg, err)
	}

	return extractSource(url.Hostname(), msg)
}

//ParseURL parses provided rawurl string into actual URL object.
func ParseURL(rawurl string) (*url.URL, error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return nil, err
	}

	return url, validate(url)
}

func extractSource(hostname string, errMsg string) (string, error) {
	tldPlusOne, err := tld.EffectiveTLDPlusOne(hostname)
	if err != nil {
		return "", fmt.Errorf(errMsg, err)
	}

	suffix, _ := tld.PublicSuffix(tldPlusOne)
	return tldPlusOne[:len(tldPlusOne)-len(suffix)-1], nil
}

func validate(subjectURL *url.URL) error {
	var err error
	validators := []func(*url.URL) error{
		validateHost,
	}

	for _, v := range validators {
		err = v(subjectURL)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateHost(url *url.URL) error {
	if url.Host == "" || url.Hostname() == "" {
		return fmt.Errorf("Unexpected URL: missing hostname")
	}

	return nil
}
