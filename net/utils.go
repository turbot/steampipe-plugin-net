package net

import (
	"net/http"
	"unicode/utf8"

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

func addRequestHeaders(request *http.Request, headers map[string]string) *http.Request {
	for header, value := range headers {
		request.Header.Add(header, value)
	}
	return request
}

type tableNetRequestRow struct {
	Url                                string
	Method                             string
	RequestBody                        string
	RequestHeaderAuthorization         string
	RequestHeaderContentType           string
	Status                             int
	StatusCode                         int
	Headers                            map[string][]string
	Body                               string
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
	Headers     map[string]string
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
