package net

import (
	"context"
	"crypto/tls"
	"time"

	"golang.org/x/exp/slices"

	"github.com/sethvargo/go-retry"
	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
)

func getQualListValues(ctx context.Context, quals map[string]*proto.QualValue, qualName string) []string {
	values := make([]string, 0)
	if quals[qualName].GetStringValue() != "" {
		values = append(values, quals[qualName].GetStringValue())
	} else if quals[qualName].GetListValue() != nil {
		for _, value := range quals[qualName].GetListValue().Values {
			str := value.GetStringValue()
			values = append(values, str)
		}
	}
	return values
}

// List all cipher suites supported by `crypto/tls` package
func cipherSuites() []*tls.CipherSuite {
	allCiphers := tls.CipherSuites()

	// also append insecure ciphers
	allCiphers = append(allCiphers, tls.InsecureCipherSuites()...)
	return allCiphers
}

// List all cipher suites supported by TLS v1.3
func cipherSuitesTLS13() []string {
	allCiphers := cipherSuites()
	var ciphersTLS13 []string

	for _, i := range allCiphers {
		if slices.Contains(i.SupportedVersions, tls.VersionTLS13) {
			ciphersTLS13 = append(ciphersTLS13, i.Name)
		}
	}
	return ciphersTLS13
}

// List all cipher suites supported by TLS v1.2
func cipherSuitesTLS12() []string {
	allCiphers := cipherSuites()
	var ciphersTLS12 []string

	for _, i := range allCiphers {
		if slices.Contains(i.SupportedVersions, tls.VersionTLS12) {
			ciphersTLS12 = append(ciphersTLS12, i.Name)
		}
	}
	return ciphersTLS12
}

// List all cipher suites supported by TLS v1.0 - TLS v1.1
func cipherSuitesUptoTLS11() []string {
	allCiphers := cipherSuites()
	var ciphersUptoTLS11 []string

	for _, i := range allCiphers {
		if slices.Contains(i.SupportedVersions, tls.VersionTLS12) || slices.Contains(i.SupportedVersions, tls.VersionTLS11) || slices.Contains(i.SupportedVersions, tls.VersionTLS10) {
			ciphersUptoTLS11 = append(ciphersUptoTLS11, i.Name)
		}
	}
	return ciphersUptoTLS11
}

// Check if given cipher suite is supported by the given protocol version
func cipherSuiteIsSupported(protocol string, cipher string) bool {
	switch protocol {
	case "TLS v1.0":
	case "TLS v1.1":
		ciphers := cipherSuitesUptoTLS11()
		return slices.Contains(ciphers, cipher)
	case "TLS v1.2":
		ciphers := cipherSuitesTLS12()
		return slices.Contains(ciphers, cipher)
	case "TLS v1.3":
		ciphers := cipherSuitesTLS13()
		return slices.Contains(ciphers, cipher)
	}
	return false
}

// Invokes the hydrate function with retryable errors and retries the function until the maximum attempts before throwing error
func retryHydrate(ctx context.Context, d *plugin.QueryData, hydrateData *plugin.HydrateData, hydrateFunc plugin.HydrateFunc) (interface{}, error) {

	// Retry configs
	retryMode := "Fibonacci"
	maxRetries := 10
	interval := 500

	// Create the backoff based on the given mode
	backoff, err := checkRetryMode(retryMode, time.Duration(interval))
	if err != nil {
		return nil, err
	}
	var hydrateResult interface{}

	err = retry.Do(ctx, retry.WithMaxRetries(uint64(maxRetries), backoff), func(ctx context.Context) error {
		hydrateResult, err = hydrateFunc(ctx, d, hydrateData)
		if err != nil {
			if shouldRetryError(err) {
				err = retry.RetryableError(err)
			}
		}
		return err
	})

	return hydrateResult, err
}

func checkRetryMode(mode string, interval time.Duration) (retry.Backoff, error) {
	switch mode {
	case "Fibonacci":
		backoff, err := retry.NewFibonacci(interval * time.Millisecond)
		if err != nil {
			return nil, err
		}
		return backoff, nil
	case "Exponential":
		backoff, err := retry.NewExponential(interval * time.Millisecond)
		if err != nil {
			return nil, err
		}
		return backoff, nil
	case "Constant":
		backoff, err := retry.NewConstant(interval * time.Millisecond)
		if err != nil {
			return nil, err
		}
		return backoff, nil
	}
	return nil, nil
}
