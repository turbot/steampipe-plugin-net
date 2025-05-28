---
title: "Steampipe Table: net_http_request - Query Network HTTP Requests using SQL"
description: "Allows users to query Network HTTP Requests, specifically details about HTTP requests and responses, providing insights into network traffic patterns and potential anomalies."
---

# Table: net_http_request - Query Network HTTP Requests using SQL

A Network HTTP Request is a service that allows you to send HTTP requests and receive HTTP responses over the network. It provides a way to interact with web services and retrieve information from web servers. Network HTTP Request helps you communicate with other services, fetch data from various sources, and interact with APIs.

## Table Usage Guide

The `net_http_request` table provides insights into HTTP requests and responses over the network. As a Network Engineer or a Developer, explore details through this table, including request method, status code, response time, and associated metadata. Utilize it to monitor network traffic, analyze performance of your web services, and troubleshoot issues related to HTTP requests and responses.

**Important Notes**
- You must specify the `url` column in the `where` clause to query this table.

## Examples

### Send a GET request to GitHub API
Explore how to evaluate the response from a specific URL, in this case, GitHub's API. This query can be used to understand the status and details of a response from a GET request to a web service, thereby aiding in API monitoring and troubleshooting.

```sql+postgres
select
  url,
  method,
  response_status_code,
  jsonb_pretty(response_body::jsonb) as response_body
from
  net_http_request
where
  url = 'https://api.github.com/users/github';
```

```sql+sqlite
select
  url,
  method,
  response_status_code,
  response_body as response_body
from
  net_http_request
where
  url = 'https://api.github.com/users/github';
```

### Send a GET request with a modified user agent
Explore the results of sending a GET request with a modified user agent. This can be useful to test how your website or application responds to different user agents.

```sql+postgres
select
  url,
  method,
  response_status_code,
  jsonb_pretty(request_headers) as request_headers,
  response_body
from
  net_http_request
where
  url = 'http://httpbin.org/user-agent'
  and request_headers = jsonb_object(
    '{user-agent, accept}',
    '{steampipe-test-user, application/json}'
  );
```

```sql+sqlite
select
  url,
  method,
  response_status_code,
  request_headers,
  response_body
from
  net_http_request
where
  url = 'http://httpbin.org/user-agent'
  and request_headers = '{"user-agent":"steampipe-test-user", "accept":"application/json"}';
```

### Send a POST request with request body and headers
This query allows you to send a POST request to a specified URL with a custom request body and headers. This can be useful for testing API endpoints, examining the response, and ensuring the server is correctly processing your request data.

```sql+postgres
select
  url,
  method,
  response_status_code,
  jsonb_pretty(request_headers) as request_headers,
  response_body
from
  net_http_request
where
  url = 'http://httpbin.org/anything'
  and method = 'POST'
  and request_body = jsonb_object(
    '{username, password}',
    '{steampipe, test_password}'
  )::text
  and request_headers = jsonb_object(
    '{content-type}',
    '{application/json}'
  );
```

```sql+sqlite
select
  url,
  method,
  response_status_code,
  request_headers,
  response_body
from
  net_http_request
where
  url = 'http://httpbin.org/anything'
  and method = 'POST'
  and request_body = '{"username":"steampipe", "password":"test_password"}'
  and request_headers = '{"content-type":"application/json"}';
```

### Send a GET request with multiple values for a request header
Discover the response details of a GET request sent to a specific URL, with multiple values for a request header. This can be useful for debugging or validating how your application handles different header configurations.

```sql+postgres
select
  url,
  method,
  response_status_code,
  jsonb_pretty(request_headers) as request_headers,
  response_body
from
  net_http_request
where
  url = 'http://httpbin.org/anything'
  and request_headers = '{
    "authorization": "Basic YWxhZGRpbjpvcGVuc2VzYW2l",
    "accept": ["application/json", "application/xml"]
  }'::jsonb;
```

```sql+sqlite
select
  url,
  method,
  response_status_code,
  request_headers,
  response_body
from
  net_http_request
where
  url = 'http://httpbin.org/anything'
  and request_headers = '{
    "authorization": "Basic YWxhZGRpbjpvcGVuc2VzYW2l",
    "accept": ["application/json", "application/xml"]
  }';
```

### Check for HTTP Strict Transport Security (HSTS) protection
Analyze the settings to understand whether a specific website, in this case, Microsoft's, has HTTP Strict Transport Security (HSTS) protection enabled. This query is useful in identifying potential security vulnerabilities related to data transmission.

```sql+postgres
select
  url,
  method,
  response_status_code,
  case
    when response_headers -> 'Strict-Transport-Security' is not null then 'enabled'
    else 'disabled'
  end as hsts_protection
from
  net_http_request
where
  url = 'http://microsoft.com';
```

```sql+sqlite
select
  url,
  method,
  response_status_code,
  case
    when json_extract(response_headers, '$.Strict-Transport-Security') is not null then 'enabled'
    else 'disabled'
  end as hsts_protection
from
  net_http_request
where
  url = 'http://microsoft.com';
```

### Send a GET request with TLS certificate verification disabled
Explore how to make a request to a site with an invalid or self-signed certificate by disabling TLS certificate verification. This is similar to using curl with the -k flag.

```sql+postgres
select
  url,
  method,
  response_status_code,
  jsonb_pretty(response_body::jsonb) as response_body
from
  net_http_request
where
  url = 'https://self-signed.example.com'
  and insecure = true;
```

```sql+sqlite
select
  url,
  method,
  response_status_code,
  response_body
from
  net_http_request
where
  url = 'https://self-signed.example.com'
  and insecure = true;
```