# Table: net_dns_reverse

Reverse DNS lookup from an IP address.

Note: An `ip_address` must be provided in all queries to this table.

## Examples

### Find host names for an IP address

```sql
select
  *
from
  net_dns_reverse
where
  ip_address = '1.1.1.1'
```
