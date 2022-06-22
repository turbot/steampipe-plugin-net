# Table: net_web_request

An HTTP request is made by a client, to a named host, which is located on a server. The aim of this table is to query the urls from Web.

Note: An `url` must be provided in all queries to this table.

## Examples

### Get properties of Turbot Organization from Github

```sql
select
  url,
  method,
  response_status_code,
  jsonb_pretty(response_body::jsonb)
from
  net_web_request
where
  url = 'https://api.github.com/users/Turbot';
```

### Make a HTTP request with no redirection

```sql
select
  url,
  method,
  response_status_code,
  jsonb_pretty(response_body::jsonb)
from
  net_web_request
where
  url = 'https://api.github.com/users/Turbot'
  and not follow_redirects;
```

### HTTP request with request context provided

`HTTP Strict Transport Security (HSTS)` is a policy mechanism that helps to protect websites against `man-in-the-middle (MITM)` attacks such as protocol downgrade attacks and cookie hijacking.

```sql
select
  url,
  method,
  response_status_code,
  jsonb_pretty(request_headers),
  response_body
from
  net_web_request
where
  url = 'http://microsoft.com'
  and request_headers = '{"authorization": "Basic YWxhZGRpbjpvcGVuc2VzYW2l", "accept": ["application/json, application/xml"]}';
```

### Check for HTTP Strict Transport Security (HSTS) protection

```sql
select
  url,
  method,
  response_status_code,
  case
    when response_headers -> 'Strict-Transport-Security' is not null then 'Enabled'
    else 'Disabled'
  end as hsts_protection
from
  net_web_request
where
  url = 'http://microsoft.com';
```
