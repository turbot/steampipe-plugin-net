# Table: net_certificate

Certificate details for a domain.

Note: A `domain` must be provided in all queries to this table.

## Examples

### Certificate information

```sql
select
  *
from
  net_certificate
where
  domain = 'steampipe.io';
```

### Time until the certificate expires

```sql
select
  domain,
  AGE(not_after, current_timestamp) as time_until_expiration
from
  net_certificate
where
  domain = 'steampipe.io';
```

### Check if the certificate is currently valid

```sql
select
  domain,
  not_before,
  not_after
from
  net_certificate
where
  domain = 'steampipe.io'
  and not_before < current_timestamp
  and not_after > current_timestamp;
```

### Check if the certificate was revoked by the CA

```sql
select
  domain,
  not_before,
  not_after
from
  net_certificate
where
  domain = 'steampipe.io'
  and is_revoked;
```
