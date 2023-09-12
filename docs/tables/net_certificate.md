# Table: net_certificate

Certificate details for a domain.

Note: A `domain` must be provided in all queries to this table.

An optional port can be included, e.g., steampipe.io:636. If not provided, the
port will default to 443.

## Examples

### Basic info

```sql
select
  *
from
  net_certificate
where
  address = 'steampipe.io';
```

### Get time until the certificate expires

```sql
select
  address,
  age(not_after, current_timestamp) as time_until_expiration
from
  net_certificate
where
  address = 'steampipe.io';
```

### Check if the certificate is currently valid

```sql
select
  address,
  not_before,
  not_after
from
  net_certificate
where
  address = 'steampipe.io',
  and not_before < current_timestamp
  and not_after > current_timestamp;
```

### Check if the certificate was revoked by the CA

```sql
select
  address,
  not_before,
  not_after
from
  net_certificate
where
  address = 'steampipe.io',
  and revoked;
```

### Check certificate revocation status with OCSP

```sql
select
  address,
  ocsp ->> 'status' as revocation_status,
  ocsp ->> 'revoked_at' as revoked_at
from
  net_certificate
where
  address = 'steampipe.io';
```

### Check if certificate using insecure algorithm (e.g., MD2, MD5, SHA1)

```sql
select
  address,
  not_before,
  not_after,
  signature_algorithm
from
  net_certificate
where
  address = 'steampipe.io',
  and signature_algorithm like any (array['%SHA1%', '%MD2%', '%MD5%']);
```

### Use a specific TLS port

```sql
select
  *
from
  net_certificate
where
  address = 'steampipe.io:464';
```
