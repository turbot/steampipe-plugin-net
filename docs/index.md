---
organization: Turbot
category: ["internet"]
icon_url: "/images/plugins/turbot/net.svg"
brand_color: "#eee"
display_name: Net
name: net
description: Network utility tables for IP addresses, DNS, certificates and more.
---

# Net

Net provides network utility tables for IP addresses, DNS, certificates and more.


## Installation

To download and install the latest net plugin:

```bash
steampipe plugin install net
```

## Credentials

None required, this plugin works with public information only.


## Connection Configuration

Connection configurations are defined using HCL in one or more Steampipe config files. Steampipe will load ALL configuration files from `~/.steampipe/config` that have a `.spc` extension. A config file may contain multiple connections.

Installing the latest net plugin will create a default connection named `net` in the `~/.steampipe/config/net.spc` file.  You must edit this connection to include your API token:

```hcl
connection "net" {
  plugin  = "net"
}
```
