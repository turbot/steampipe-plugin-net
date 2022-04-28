# Table: net_dns_record

DNS records associated with a given domain.

The default DNS server used for all requests is the Google global public server, 8.8.8.8. This default can be overriden in 2 ways:

- Update the `dns_server` configuration argument.
- Specify `dns_server` in the query, which overrides the default and `dns_server` configuration argument. For instance, to use Cloudflare's global public server instead:
  ```sql
  select
    *
  from
    net_dns_record
  where
    domain = 'steampipe.io'
    and dns_server = '1.1.1.1:53';
  ```

Note: A `domain` must be provided in all queries to this table.

## Examples

### DNS records for a domain

```sql
select
  *
from
  net_dns_record
where
  domain = 'steampipe.io'
```

### List TXT records for a domain

```sql
select
  value,
  ttl
from
  net_dns_record
where
  domain = 'github.com'
  and type = 'TXT'
```

### Mail server records for a domain in priority order

```sql
select
  target,
  priority,
  ttl
from
  net_dns_record
where
  domain = 'turbot.com'
  and type = 'MX'
order by
  priority
```
