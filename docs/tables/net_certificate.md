# Table: net_certificate

Certificate details for a domain.

Note: An `address` of the format address:port (e.g., steampipe.io:443) must be provided.

## Examples

### Basic info

```sql
select
  *
from
  net_certificate
where
  address = 'steampipe.io:443';
```

### Get time until the certificate expires

```sql
select
  address,
  age(not_after, current_timestamp) as time_until_expiration
from
  net_certificate
where
  address = 'steampipe.io:443';
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
  address = 'steampipe.io:443'
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
  address = 'steampipe.io:443'
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
  address = 'steampipe.io:443';
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
  address = 'steampipe.io:443'
  and signature_algorithm like any (array['%SHA1%', '%MD2%', '%MD5%']);
```
