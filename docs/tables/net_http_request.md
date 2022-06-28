# Table: net_http_request

Send an HTTP request to a host on a server in order to receive a response. A request body and headers can be provided with each request. Currently, the GET and POST methods are supported.

Note: A `url` must be provided in all queries to this table.

## Examples

### Send a GET request to GitHub API

```sql
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

### Send a GET request with a modified user agent

```sql
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

### Send a POST request with request body and headers

```sql
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

### Send a GET request with multiple values for a request header

```sql
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

### Check for HTTP Strict Transport Security (HSTS) protection

```sql
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
