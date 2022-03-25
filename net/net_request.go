package net

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v3/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v3/plugin/transform"
)

func tableNetRequest() *plugin.Table {
	return &plugin.Table{
		Name: "net_request",
		List: &plugin.ListConfig{
			ParentHydrate: listBaseRequestAttributes,
			Hydrate:       listRequestResponses,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "url", Require: plugin.Required},
				{Name: "method", Require: plugin.Optional},
				{Name: "request_header_authorization", Require: plugin.Optional},
				{Name: "request_header_content_type", Require: plugin.Optional},
				{Name: "request_body", Require: plugin.Optional},
			},
		},
		Columns: []*plugin.Column{
			{Name: "url", Transform: transform.FromField("Url"), Type: proto.ColumnType_STRING},
			{Name: "method", Type: proto.ColumnType_STRING},
			{Name: "request_body", Type: proto.ColumnType_STRING},
			{Name: "request_header_authorization", Type: proto.ColumnType_STRING},
			{Name: "request_header_content_type", Type: proto.ColumnType_STRING},
			{Name: "status_code", Type: proto.ColumnType_INT},
			{Name: "body", Type: proto.ColumnType_STRING},
			{Name: "error", Type: proto.ColumnType_STRING},
			{Name: "headers", Type: proto.ColumnType_JSON},
			{Name: "header_content_security_policy", Type: proto.ColumnType_STRING},
			{Name: "header_content_type", Type: proto.ColumnType_STRING},
			{Name: "header_permissions_policy", Type: proto.ColumnType_STRING},
			{Name: "header_referrer_policy", Type: proto.ColumnType_STRING},
			{Name: "header_x_frame_options", Type: proto.ColumnType_STRING},
			{Name: "header_x_xss_protection", Type: proto.ColumnType_STRING},
		},
	}
}

func listBaseRequestAttributes(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	logger := plugin.Logger(ctx)

	var methods []string
	var requestBody string
	headers := make(map[string]string)

	queryCols := d.KeyColumnQuals

	urls := getQuals(queryCols["url"])
	if queryCols["method"] != nil {
		methods = getQuals(queryCols["method"])
	} else {
		methods = []string{"GET"}
	}

	if authHeader, present := getAuthHeaderQuals(queryCols["request_header_authorization"]); present {
		headers["Authorization"] = authHeader
	}

	if contentTypeHeader, present := getAuthHeaderQuals(queryCols["request_header_content_type"]); present {
		headers["Content-Type"] = contentTypeHeader
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

	url := baseRequestAttribute.Url
	methods := baseRequestAttribute.Methods
	headers := baseRequestAttribute.Headers
	requestBody := baseRequestAttribute.RequestBody
	client := &http.Client{}

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

		req = addRequestHeaders(req, headers)

		// Make request
		res, requestErr = client.Do(req)

		// Read response
		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(res.Body)
		if err != nil {
			logger.Error("listRequestResponses buffer reading error", "url", url, "request method", req.Method, "error", err)
			return nil, err
		}

		// Close response reading
		res.Body.Close()
		body := buf.String()

		// Generate table row item
		item := tableNetRequestRow{
			Url:                                url,
			Method:                             method,
			RequestBody:                        requestBody,
			RequestHeaderContentType:           headers["Content-Type"],
			RequestHeaderAuthorization:         headers["Authorization"],
			StatusCode:                         res.StatusCode,
			Headers:                            res.Header,
			Body:                               body,
			HeaderContentSecurityPolicy:        res.Header.Get("Content-Security-Policy"),
			HeaderContentType:                  res.Header.Get("Content-Type"),
			HeaderCrossSiteScriptingProtection: res.Header.Get("X-XSS-Protection"),
			HeaderPermissionsPolicy:            res.Header.Get("Permissions-Policy"),
			HeaderReferrerPolicy:               res.Header.Get("Referrer-Policy"),
			HeaderXFrameOptions:                res.Header.Get("X-Frame-Options"),
		}
		if requestErr != nil {
			logger.Error("listRequestResponses do request error", "url", url, "request method", req.Method, "error", requestErr.Error())
			item.Error = requestErr.Error()
		}

		d.StreamListItem(ctx, item)
	}

	return nil, nil
}