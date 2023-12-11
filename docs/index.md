---
organization: Turbot
category: ["internet"]
icon_url: "/images/plugins/turbot/net.svg"
brand_color: "#005A9C"
display_name: Net
name: net
description: Steampipe plugin for querying DNS records, certificates and other network information.
og_description: Query networking information with SQL! Zero ETL CLI. No DB required. 
og_image: "/images/plugins/turbot/net-social-graphic.png"
engines: ["steampipe", "sqlite", "postgres", "export"]
---

# Net + Steampipe

[Steampipe](https://steampipe.io) is a CLI to instantly query cloud APIs using SQL.

The net plugin is a set of utility tables for steampipe to query attributes of X.509 certificates associated with websites, DNS records, and connectivity to specific network socket addresses.

For example:

```sql
select
  issuer, 
  not_after as exp_date 
from 
  net_certificate
where
  domain = 'steampipe.io';
```

```sh
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
| Radius | Steampipe limits searches to specific resources based on the provided `Quals` e.g. `domain` for certificates and DNS queries and `address` for network connection information |
| Resolution | n/a |

### Configuration

No configuration is needed. Installing the latest net plugin will create a config file (`~/.steampipe/config/net.spc`) with a single connection named `net`:

```hcl
connection "net" {
  plugin = "net"
}
```

## Get involved

- GitHub: https://github.com/turbot/steampipe-plugin-net
- Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
