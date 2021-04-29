# Table: net_connection

Test network connectivity to using a network protocol (e.g. TCP) and address / port (e.g. steampipe.io:443).

Note: The `protocol` and `address` columns must be provided in all queries to this table.

## Examples

### Test a TCP connection to steampipe.io on port 443

```sql
select
  *
from
  net_connection
where
  protocol = 'tcp'
  and address = 'steampipe.io:443'
```

### Test a UDP connection to DNS server 1.1.1.1 on port 53

```sql
select
  *
from
  net_connection
where
  protocol = 'udp'
  and address = '1.1.1.1:53'
```

### Test if SSH is open on server 68.183.153.44

```sql
select
  *
from
  net_connection
where
  protocol = 'tcp'
  and address = '68.183.153.44:ssh'
```

### Test if RDP is open on server 65.2.9.152

```sql
select
  *
from
  net_connection
where
  protocol = 'tcp'
  and address = '65.2.9.152:3389'
```
