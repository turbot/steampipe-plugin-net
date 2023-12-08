---
title: "Steampipe Table: net_connection - Query Network Connection using SQL"
description: "Allows users to query Network Connections, specifically providing insights into various network connections, their statuses, and related details."
---

# Table: net_connection - Query Network Connection using SQL

A Network Connection is a link between two or more nodes in a network. It enables the transfer of data between these nodes, which can include computers, servers, or other network-enabled devices. The status and details of these connections are crucial for network management and troubleshooting.

## Table Usage Guide

The `net_connection` table provides insights into various network connections. As a Network Engineer or IT Administrator, you can explore connection-specific details through this table, including statuses, types, and associated metadata. Utilize it to monitor and manage network connections, ensure optimal data transfer, and troubleshoot any connection issues.

## Examples

### Test a TCP connection (the default protocol) to steampipe.io on port 443
Analyze the status of a TCP connection to a specific website and port. This can be useful for troubleshooting network connectivity issues or verifying that a service is reachable and responding as expected.

```sql+postgres
select
  *
from
  net_connection
where
  address = 'steampipe.io:443';
```

```sql+sqlite
select
  *
from
  net_connection
where
  address = 'steampipe.io:443';
```

### Test if SSH is open on server 68.183.153.44
The query allows you to assess if a specific server has an open SSH connection. This is useful for identifying potential security vulnerabilities or for troubleshooting connectivity issues.

```sql+postgres
select
  *
from
  net_connection
where
  address = '68.183.153.44:ssh';
```

```sql+sqlite
select
  *
from
  net_connection
where
  address = '68.183.153.44:ssh';
```

### Test a UDP connection to DNS server 1.1.1.1 on port 53
Explore whether a UDP connection to a DNS server on a specific port is active. This is useful to troubleshoot network connectivity issues or validate network configurations.

```sql+postgres
select
  *
from
  net_connection
where
  protocol = 'udp'
  and address = '1.1.1.1:53';
```

```sql+sqlite
select
  *
from
  net_connection
where
  protocol = 'udp'
  and address = '1.1.1.1:53';
```

### Test if RDP is open on server 65.2.9.152
Explore whether the Remote Desktop Protocol (RDP) is open on a specific server to ensure secure connections and prevent unauthorized access. This is particularly useful in managing network security and maintaining control over remote access to your systems.

```sql+postgres
select
  *
from
  net_connection
where
  protocol = 'tcp'
  and address = '65.2.9.152:3389';
```

```sql+sqlite
select
  *
from
  net_connection
where
  protocol = 'tcp'
  and address = '65.2.9.152:3389';
```