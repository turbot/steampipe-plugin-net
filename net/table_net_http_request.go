package net

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

func tableNetHTTPRequest() *plugin.Table {
	return &plugin.Table{
		Name:        "net_http_request",
		Description: "An HTTP request is made by a client, to a named host, which is located on a server.",
		List: &plugin.ListConfig{
			ParentHydrate: listBaseRequestAttributes,
			Hydrate:       listRequestResponses,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "url", Require: plugin.Required},
				{Name: "method", Require: plugin.Optional, CacheMatch: "exact"},
				{Name: "follow_redirects", Require: plugin.Optional, Operators: []string{"=", "<>"}, CacheMatch: "exact"},
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
			{Name: "response_error", Type: proto.ColumnType_STRING, Description: "Represents an error or failure, either from a non-successful HTTP status, an error while executing the request, or some other failure which occurred during the parsing of the response.", Transform: transform.FromField("Error")},
			{Name: "response_headers", Type: proto.ColumnType_JSON, Description: "A map of response headers used by web applications to configure security defenses in web browsers."},
		},
	}
}

func listBaseRequestAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var methods []string
	var requestBody string
	headers := make(map[string]interface{})

	queryCols := d.EqualsQuals

	urls := getQuals(queryCols["url"])
	logger.Debug("listBaseRequestAttributes", "urls", urls)
	logger.Debug("listBaseRequestAttributes", "query cols", queryCols)
	if queryCols["method"] != nil {
		methods = getQuals(queryCols["method"])
	} else {
		methods = []string{"GET"}
	}

	requestHeadersString := queryCols["request_headers"].GetJsonbValue()
	logger.Debug("listBaseRequestAttributes", "request headers", requestHeadersString)

	if requestHeadersString != "" {
		err := json.Unmarshal([]byte(requestHeadersString), &headers)
		if err != nil {
			plugin.Logger(ctx).Error("net_http_request.listBaseRequestAttributes", "unmarshal_error", err)
			return nil, fmt.Errorf("failed to unmarshal request headers: %v", err)
		}
	}

	for k, v := range headers {
		logger.Debug("listBaseRequestAttributes", "header", k, v)
	}

	if requestBodyData, present := getAuthHeaderQuals(queryCols["request_body"]); present {
		requestBody = requestBodyData
	}

	logger.Debug("listBaseRequestAttributes", "urls", urls)
	logger.Debug("listBaseRequestAttributes", "methods", methods)
	logger.Debug("listBaseRequestAttributes", "headers", headers)

	for _, url := range urls {
		d.StreamListItem(ctx, baseRequestAttributes{url, methods, requestBody, headers})
	}

	return nil, nil
}

func listRequestResponses(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)
	baseRequestAttribute := h.Item.(baseRequestAttributes)

	logger.Debug("listRequestResponses", "attributes", baseRequestAttribute)

	url := baseRequestAttribute.Url
	methods := baseRequestAttribute.Methods
	headers := baseRequestAttribute.Headers
	requestBody := baseRequestAttribute.RequestBody
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
			logger.Warn("listRequestResponses", "unsupported method", method)
			continue
		}

		req = addRequestHeaders(req, headers)
		logger.Debug("listRequestResponses", "request", req)

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

			item.ResponseStatusCode = res.StatusCode
			item.ResponseHeaders = res.Header
			item.ResponseBody = body
		}

		// Generate table row item
		d.StreamListItem(ctx, item)
	}

	return nil, nil
}
