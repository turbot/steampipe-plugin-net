---
organization: Turbot
category: ["internet"]
icon_url: "/images/plugins/turbot/net.svg"
brand_color: "#005A9C"
display_name: Net
name: net
description: Steampipe plugin for querying DNS records, certificates and other network information.
og_description: Query networking information with SQL! Open source CLI. No DB required. 
og_image: "/images/plugins/turbot/net-social-graphic.png"
---

# Net + Steampipe

[Steampipe](https://steampipe.io) is an open source CLI to instantly query cloud APIs using SQL.

The net plugin is a set of utility tables for steampipe to query attributes of X.509 certificates associated with websites, DNS records and connectivity to specific network socket addresses.

For example:

```sql
select
  issuer, 
  not_before as exp_date 
from 
  net_certificate
where
  domain = 'steampipe.io';
```

```
+----------------------------+---------------------+
| issuer                     | exp_date            |
+----------------------------+---------------------+
| CN=R3,O=Let's Encrypt,C=US | 2021-02-24 03:02:15 |
+----------------------------+---------------------+
```

## Documentation

- **[Table definitions & examples →](/plugins/turbot/net/tables)**

## Get started

### Install

Download and install the latest Steampipe Net plugin:

```bash
steampipe plugin install net
```

### Credentials

| Item | Description |
| - | - |
| Credentials | No creds required |
| Permissions | n/a |
| Radius | Steampipe limits searches to specific resources based on the provided `Quals` e.g. `domain` for certificates and dns queries and `address` for network connection information |
| Resolution | n/a |

### Configuration

No configuration is needed. Installing the latest net plugin will create a config file (`~/.steampipe/config/net.spc`) with a single connection named `net`:

```hcl
connection "net" {
  plugin = "net"
}
```

## Get involved

* Open source: https://github.com/turbot/steampipe-plugin-net
* Community: [Slack Channel](https://steampipe.io/community/join)