# Table: net_tls_connection

Test TLS connection to the given network address (e.g. steampipe.io:443) by initiating a TLS handshake. This table checks connection for all possible TLS protocol-cipher combination and returns the combinations for which a TLS connection could be established.

Note: An `address` of the format domain:port (e.g. steampipe.io:443) must be provided.

You can also provide a protocol version and cipher suite to verify the TLS connection as per requirement. For example:

```sql
select * from net_tls_connection where address = 'steampipe.io:443' and version = 'TLS v1.3' and cipher_suite_name = 'TLS_AES_128_GCM_SHA256';
```

Note:

- SSL protocols (e.g. SSL v3 and SSL v2) are not supported by this table.
- This table supports a limited set of cipher suites (supported by https://pkg.go.dev/crypto/tls package).

## Examples

### Check all the supported protocols and cipher suites for which a TLS connection could be established

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
  and handshake_completed;
```

### Check TLS handshake with a certain protocol and cipher suite

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

### Check if server allows to connect with an insecure cipher suite

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
  and cipher_suite_name in ('TLS_RSA_WITH_RC4_128_SHA', 'TLS_RSA_WITH_3DES_EDE_CBC_SHA', 'TLS_RSA_WITH_AES_128_CBC_SHA256', 'TLS_ECDHE_ECDSA_WITH_RC4_128_SHA', 'TLS_ECDHE_RSA_WITH_RC4_128_SHA', 'TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA', 'TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256', 'TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256')
  and handshake_completed;
```
