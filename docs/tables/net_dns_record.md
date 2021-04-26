# Table: net_dns_record

DNS records associated with a given domain.

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
