---
title: "Steampipe Table: net_dns_record - Query DNS Records using SQL"
description: "Allows users to query DNS Records, specifically the details of each DNS record, providing insights into the DNS configuration and potential issues."
---

# Table: net_dns_record - Query DNS Records using SQL

 DNS service is a scalable, reliable, and managed Domain Name System (DNS) service that provides a high-performance, global footprint for your public-facing internet resources and resolves DNS zones. It enables you to distribute traffic to your endpoints and ensures high availability and failover, thus improving the performance of your web applications.

## Table Usage Guide

The `net_dns_record` table provides insights into DNS Records. As a network administrator, you can explore DNS record-specific details through this table, including record types, domain names, and associated metadata. Utilize it to uncover information about DNS records, such as those with misconfigured settings, the association between domain names and IP addresses, and the verification of DNS record settings.

**Important Notes**

The default DNS server used for all requests is the Google global public server, 8.8.8.8. This default can be overriden in 2 ways:

- Update the `dns_server` configuration argument.
- Specify `dns_server` in the query, which overrides the default and `dns_server` configuration argument. For instance, to use Cloudflare's global public server instead:

```sql+postgres
select
  *
from
  net_dns_record
where
  domain = 'steampipe.io'
  and dns_server = '1.1.1.1:53';
```

```sql+sqlite
select
  *
from
  net_dns_record
where
  domain = 'steampipe.io'
  and dns_server = '1.1.1.1:53';
```

- A `domain` must be provided in all queries to this table.

## Examples

### DNS records for a domain
Explore DNS records associated with a specific domain to understand its configuration and structure. This could be beneficial for troubleshooting or auditing purposes.

```sql+postgres
select
  *
from
  net_dns_record
where
  domain = 'steampipe.io';
```

```sql+sqlite
select
  *
from
  net_dns_record
where
  domain = 'steampipe.io';
```

### List TXT records for a domain
Explore the text records for a specific domain to understand its associated data and time-to-live values. This could be useful for verifying domain ownership or understanding security settings.

```sql+postgres
select
  value,
  ttl
from
  net_dns_record
where
  domain = 'github.com'
  and type = 'TXT';
```

```sql+sqlite
select
  value,
  ttl
from
  net_dns_record
where
  domain = 'github.com'
  and type = 'TXT';
```

### Mail server records for a domain in priority order
Explore the priority order of mail servers for a specific domain. This is beneficial for understanding the order in which email will be delivered or rerouted if the primary server is not available.

```sql+postgres
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
  priority;
```

```sql+sqlite
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
  priority;
```