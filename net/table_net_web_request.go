package net

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableNetWebRequest() *plugin.Table {
	return &plugin.Table{
		Name: "net_web_request",
		List: &plugin.ListConfig{
			ParentHydrate: listBaseRequestAttributes,
			Hydrate:       listRequestResponses,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "url", Require: plugin.Required},
				{Name: "method", Require: plugin.Optional},
				{Name: "follow_redirects", Require: plugin.Optional, Operators: []string{"=", "<>"}},
				{Name: "request_headers", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "request_body", Require: plugin.Optional, CacheMatch: "exact"},
			},
		},
		Columns: []*plugin.Column{
			{Name: "url", Transform: transform.FromField("Url"), Type: proto.ColumnType_STRING, Description: "URL of the site."},
			{Name: "method", Type: proto.ColumnType_STRING, Description: "Specifies the HTTP method (GET, POST)."},
			{Name: "follow_redirects", Type: proto.ColumnType_BOOL, Description: "If true, the requests will follow the redirects."},
			{Name: "request_body", Type: proto.ColumnType_STRING, Description: "The request's body."},
			{Name: "request_headers", Type: proto.ColumnType_JSON, Transform: transform.FromQual("request_headers"), Description: "A map of headers passed in the request."},
			{Name: "response_status_code", Type: proto.ColumnType_INT, Description: "HTTP status code is a server response to a browser's request."},
			{Name: "response_body", Type: proto.ColumnType_STRING, Description: "Represents the response body."},
			// TODO: Does it need response_? What is this?
			{Name: "response_error", Type: proto.ColumnType_STRING, Description: "Represents an error or failure, either from a non-successful HTTP status, an error while executing the request, or some other failure which occurred during the parsing of the response.", Transform: transform.FromField("Error")},
			{Name: "response_headers", Type: proto.ColumnType_JSON, Description: "A map of response headers used by web applications to configure security defenses in web browsers."},
		},
	}
}

func listBaseRequestAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var methods []string
	var requestBody string
	headers := make(map[string]string)

	logger.Info("listBaseRequestAttributes", "Headers", headers)

	queryCols := d.KeyColumnQuals

	urls := getQuals(queryCols["url"])
	logger.Info("listBaseRequestAttributes", "URLs", urls)
	logger.Info("listBaseRequestAttributes", "Query cols", queryCols)
	if queryCols["method"] != nil {
		methods = getQuals(queryCols["method"])
	} else {
		methods = []string{"GET"}
	}

	requestHeadersString := queryCols["request_headers"].GetJsonbValue()
	logger.Info("listBaseRequestAttributes", "Headers String", requestHeadersString)

	// TODO: How to handle headers with same key and different values? Use comma delimited?
	if requestHeadersString != "" {
		json.Unmarshal([]byte(requestHeadersString), &headers)
	}

	for k, v := range headers {
		logger.Info("listBaseRequestAttributes", "Header", k, v)
	}

	if requestBodyData, present := getAuthHeaderQuals(queryCols["request_body"]); present {
		requestBody = requestBodyData
	}

	logger.Info("listBaseRequestAttributes", "URLs", urls)
	logger.Info("listBaseRequestAttributes", "Methods", methods)
	logger.Info("listBaseRequestAttributes", "Headers", headers)

	for _, url := range urls {
		d.StreamListItem(ctx, baseRequestAttributes{url, methods, requestBody, headers})
	}

	return nil, nil
}

func listRequestResponses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	baseRequestAttribute := h.Item.(baseRequestAttributes)

	logger.Info("listRequestResponses", "Attributes", baseRequestAttribute)

	url := baseRequestAttribute.Url
	methods := baseRequestAttribute.Methods
	headers := baseRequestAttribute.Headers
	requestBody := baseRequestAttribute.RequestBody

	// TODO: Should this be an argument? Default to false (secure by default)?
	//tr := &http.Transport{
	//	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	//}
	//client := &http.Client{Transport: tr}
	client := &http.Client{}

	// Set true to follow the redirects
	// Default set to true
	// If passed using follow_redirects, override the default
	followRedirects := true
	if d.Quals["follow_redirects"] != nil {
		for _, q := range d.Quals["follow_redirects"].Quals {
			switch q.Operator {
			case "<>":
				followRedirects = false
			case "=":
				followRedirects = true
			}
		}
	}

	if !followRedirects {
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	// Execute the request for each type of method per url
	for _, method := range methods {
		var req *http.Request
		var res *http.Response
		var err error
		var requestErr error

		switch method {
		case "GET":
			req, err = http.NewRequest("GET", url, nil)
			if err != nil {
				logger.Error("listRequestResponses", "url", url, "New GET Request error", err)
				return nil, err
			}
		case "POST":
			requestBodyReader := strings.NewReader(requestBody)
			req, err = http.NewRequest("POST", url, requestBodyReader)
			if err != nil {
				logger.Error("listRequestResponses", "url", url, "New Post Request error", err)
				return nil, err
			}
		default:
			logger.Warn("Unsupported method", method)
			continue
		}

		// req.Proto = "HTTP/2"

		req = addRequestHeaders(req, headers)
		logger.Info("Request:", req)

		item := tableNetWebRequestRow{
			Url:             url,
			Method:          method,
			RequestBody:     requestBody,
			FollowRedirects: followRedirects,
		}

		// Make request
		res, requestErr = client.Do(req)
		if requestErr != nil {
			logger.Error("listRequestResponses do request error", "url", url, "request method", req.Method, "error", requestErr.Error())
			item.Error = requestErr.Error()
		}

		if requestErr == nil {
			// Read response
			buf := new(bytes.Buffer)
			_, err = buf.ReadFrom(res.Body)
			if err != nil {
				logger.Error("listRequestResponses buffer reading error", "url", url, "request method", req.Method, "error", err)
				return nil, err
			}

			// Close response reading
			res.Body.Close()
			body := removeInvalidUTF8Char(buf.String())

			queryCols := d.KeyColumnQuals
			requestHeadersString := queryCols["request_headers"].GetStringValue()
			logger.Info("listRequestResponses", "Headers String", requestHeadersString)

			item.ResponseStatusCode = res.StatusCode
			item.ResponseHeaders = res.Header
			item.ResponseBody = body
		}

		// // TODO: Can we show the full redirect res chain?
		// // TODO: What cert info do we get?
		// // Generate table row item
		d.StreamListItem(ctx, item)
	}

	return nil, nil
}
