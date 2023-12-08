---
title: "Steampipe Table: net_tls_connection - Query Network TLS Connections using SQL"
description: "Allows users to query Network TLS Connections, specifically details about SSL/TLS connections established by the system, providing insights into connection parameters and potential security issues."
---

# Table: net_tls_connection - Query Network TLS Connections using SQL

A Network TLS Connection is a secure connection established between two systems using the Transport Layer Security (TLS) protocol. This protocol provides privacy and data integrity between applications communicating over a network. Information about these connections can be useful in understanding network traffic patterns and identifying potential security vulnerabilities.

## Table Usage Guide

The `net_tls_connection` table provides insights into the SSL/TLS connections established by your system. As a network administrator or security analyst, explore connection-specific details through this table, including connection parameters, encryption algorithms used, and certificate details. Utilize it to uncover information about connections, such as those using weak encryption algorithms, expired certificates, and potential security vulnerabilities.

## Examples

### List all supported protocols and cipher suites for which a TLS connection could be established
Explore which protocols and cipher suites can successfully establish a TLS connection to a specific address. This can be useful to ensure secure and compatible connections in your network communication.

```sql+postgres
select
  address,
  version,
  cipher_suite_name,
  handshake_completed
from
  net_tls_connection
where
  address = 'steampipe.io:443'
  and handshake_completed;
```

```sql+sqlite
select
  address,
  version,
  cipher_suite_name,
  handshake_completed
from
  net_tls_connection
where
  address = 'steampipe.io:443'
  and handshake_completed;
```

### Check TLS handshake with a certain protocol and cipher suite
Identify instances where a specific protocol and cipher suite have successfully completed a TLS handshake with a particular server. This can be useful to ensure secure communication and confirm the server's compatibility with desired security standards.

```sql+postgres
select
  address,
  version,
  cipher_suite_name,
  handshake_completed
from
  net_tls_connection
where
  address = 'steampipe.io:443'
  and version = 'TLS v1.3'
  and cipher_suite_name = 'TLS_AES_128_GCM_SHA256';
```

```sql+sqlite
The given PostgreSQL query does not use any PostgreSQL-specific functions or data types, nor does it involve any JSON manipulation or join operations. Therefore, it can be directly used in SQLite without any modification. Here is the equivalent SQLite query:

```sql
select
  address,
  version,
  cipher_suite_name,
  handshake_completed
from
  net_tls_connection
where
  address = 'steampipe.io:443'
  and version = 'TLS v1.3'
  and cipher_suite_name = 'TLS_AES_128_GCM_SHA256';
```
```

### Check if a server allows connections with an insecure cipher suite
Determine if a server is susceptible to security risks by identifying connections that use potentially insecure cipher suites. This can be useful for enhancing security measures and preventing potential cyber threats.

```sql+postgres
select
  address,
  version,
  cipher_suite_name,
  handshake_completed
from
  net_tls_connection
where
  address = 'steampipe.io:443'
  and cipher_suite_name in ('TLS_RSA_WITH_RC4_128_SHA', 'TLS_RSA_WITH_3DES_EDE_CBC_SHA', 'TLS_RSA_WITH_AES_128_CBC_SHA256', 'TLS_ECDHE_ECDSA_WITH_RC4_128_SHA', 'TLS_ECDHE_RSA_WITH_RC4_128_SHA', 'TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA', 'TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256', 'TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256')
  and handshake_completed;
```

```sql+sqlite
select
  address,
  version,
  cipher_suite_name,
  handshake_completed
from
  net_tls_connection
where
  address = 'steampipe.io:443'
  and cipher_suite_name in ('TLS_RSA_WITH_RC4_128_SHA', 'TLS_RSA_WITH_3DES_EDE_CBC_SHA', 'TLS_RSA_WITH_AES_128_CBC_SHA256', 'TLS_ECDHE_ECDSA_WITH_RC4_128_SHA', 'TLS_ECDHE_RSA_WITH_RC4_128_SHA', 'TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA', 'TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256', 'TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256')
  and handshake_completed = 1;
```