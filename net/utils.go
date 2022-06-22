package net

import (
	"context"
	"crypto/tls"
	"net/http"
	"unicode/utf8"

	"golang.org/x/exp/slices"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
)

func getQuals(qualValue *proto.QualValue) []string {
	var data []string
	// Check for nil
	if qualValue == nil {
		return data
	}
	if qualList := qualValue.GetListValue(); qualList != nil {
		for _, q := range qualList.Values {
			data = append(data, q.GetStringValue())
		}
	} else {
		data = append(data, qualValue.GetStringValue())
	}
	return data
}

func getAuthHeaderQuals(qualValue *proto.QualValue) (authHeader string, present bool) {
	if qualValue == nil {
		return "", false
	}
	if qualList := qualValue.GetListValue(); qualList != nil {
		for _, q := range qualList.Values {
			return q.GetStringValue(), true
		}
	}
	return qualValue.GetStringValue(), true
}

func addRequestHeaders(request *http.Request, headers map[string]interface{}) *http.Request {
	for header, value := range headers {
		content, isArray := value.([]interface{})
		if isArray {
			for _, i := range content {
				request.Header.Add(header, i.(string))
			}
		} else {
			request.Header.Add(header, value.(string))
		}
	}
	return request
}

type tableNetWebRequestRow struct {
	Url             string
	Method          string
	RequestBody     string
	RequestHeaders  string
	FollowRedirects bool
	//RequestHeaderAuthorization         string
	//RequestHeaderContentType           string
	Status                             int
	ResponseStatusCode                 int
	ResponseHeaders                    map[string][]string
	ResponseBody                       string
	Error                              string
	HeaderContentSecurityPolicy        string
	HeaderStrictTransportSecurity      string
	HeaderContentType                  string
	HeaderCrossSiteScriptingProtection string
	HeaderPermissionsPolicy            string
	HeaderReferrerPolicy               string
	HeaderXFrameOptions                string
	HeaderXContentTypeOptions          string
}

type baseRequestAttributes struct {
	Url         string
	Methods     []string
	RequestBody string
	Headers     map[string]interface{}
}

func removeInvalidUTF8Char(str string) string {
	if !utf8.ValidString(str) {
		v := make([]rune, 0, len(str))
		for i, r := range str {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(str[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		str = string(v)
	}
	return str
}

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
