---
title: "Steampipe Table: net_certificate - Query Net Certificates using SQL"
description: "Allows users to query Net Certificates, providing details about the certificate's validity, issuer, subject, and other related information."
---

# Table: net_certificate - Query Net Certificates using SQL

A Net Certificate is a digital document that verifies a server's details. When a browser initiates a connection with a secure website, the web server sends its public key certificate for the browser to check. This document contains the server's public key, the certificate's validity dates, and an identifier for the certificate authority (CA) that issued the certificate.

## Table Usage Guide

The `net_certificate` table provides insights into Net Certificates. As a Security Analyst, explore certificate-specific details through this table, including issuer, subject, validity, and associated metadata. Utilize it to uncover information about certificates, such as their validity status, issuer details, and the verification of the certificate authority.

## Examples

### Basic info
Analyze the settings to understand the security certificates associated with a specific web address, such as 'steampipe.io:443'. This can be useful for assessing the security status and identifying any potential issues or vulnerabilities.

```sql+postgres
select
  *
from
  net_certificate
where
  address = 'steampipe.io:443';
```

```sql+sqlite
select
  *
from
  net_certificate
where
  address = 'steampipe.io:443';
```

### Get time until the certificate expires
Determine the remaining validity period of a specific certificate. This query is useful in monitoring and ensuring that the certificate does not expire unexpectedly, thereby preventing potential service interruptions.

```sql+postgres
select
  address,
  age(not_after, current_timestamp) as time_until_expiration
from
  net_certificate
where
  address = 'steampipe.io:443';
```

```sql+sqlite
select
  address,
  julianday(not_after) - julianday(current_timestamp) as time_until_expiration
from
  net_certificate
where
  address = 'steampipe.io:443';
```

### Check if the certificate is currently valid
Explore which security certificates are currently valid by assessing their validity periods. This is particularly useful to ensure your connections to certain addresses, like 'steampipe.io:443', are secure and up-to-date.

```sql+postgres
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

```sql+sqlite
select
  address,
  not_before,
  not_after
from
  net_certificate
where
  address = 'steampipe.io:443'
  and not_before < datetime('now')
  and not_after > datetime('now');
```

### Check if the certificate was revoked by the CA
Determine if a specific website's security certificate has been revoked by the certificate authority. This is useful for understanding the security status of your web connections.

```sql+postgres
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

```sql+sqlite
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
Determine the revocation status of a certificate using Online Certificate Status Protocol (OCSP). This can help in identifying if a certificate has been revoked, and if so, when it happened, which is crucial for maintaining secure online connections.

```sql+postgres
select
  address,
  ocsp ->> 'status' as revocation_status,
  ocsp ->> 'revoked_at' as revoked_at
from
  net_certificate
where
  address = 'steampipe.io:443';
```

```sql+sqlite
select
  address,
  json_extract(ocsp, '$.status') as revocation_status,
  json_extract(ocsp, '$.revoked_at') as revoked_at
from
  net_certificate
where
  address = 'steampipe.io:443';
```

### Check if certificate using insecure algorithm (e.g., MD2, MD5, SHA1)
Explore which digital certificates are using insecure algorithms, such as MD2, MD5, or SHA1. This query is beneficial for identifying potential security risks associated with outdated or weak cryptographic algorithms.

```sql+postgres
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

```sql+sqlite
select
  address,
  not_before,
  not_after,
  signature_algorithm
from
  net_certificate
where
  address = 'steampipe.io:443'
  and (signature_algorithm like '%SHA1%' or signature_algorithm like '%MD2%' or signature_algorithm like '%MD5%');
```