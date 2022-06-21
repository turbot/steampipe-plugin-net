## v0.4.0 [2022-06-09]

_What's new?_

- New tables added: ([#35](https://github.com/turbot/steampipe-plugin-net/pull/35))
  - [net_tls_connection](https://hub.steampipe.io/plugins/turbot/net/tables/net_tls_connection)

_Enhancements_

- Added `tag` column to `net_dns_record` table. ([#33](https://github.com/turbot/steampipe-plugin-net/pull/33))
- Added `crl_distribution_points`, `issuer_name`, `ocsp`, `ocsp_servers`, `public_key_length`, `revoked`, and `transparent` columns to `net_certificate` table. ([#33](https://github.com/turbot/steampipe-plugin-net/pull/33))

## v0.3.0 [2022-04-28]

_Enhancements_

- Added columns `dns_server`, `expire`, `minimum`, `refresh`, `retry`, `serial` to `net_dns_record` table. ([#28](https://github.com/turbot/steampipe-plugin-net/pull/28))
- Updated  `net_dns_record` table to use Google's global public DNS instead of Cloudflare's for faster results. ([#28](https://github.com/turbot/steampipe-plugin-net/pull/28))
- Recompiled plugin with miekg/dns v1.1.47. ([#28](https://github.com/turbot/steampipe-plugin-net/pull/28))

_Bug fixes_

- Fixed `net_dns_record` table not returning correct results for consecutive queries when using the `type` list key column. ([#28](https://github.com/turbot/steampipe-plugin-net/pull/28))

## v0.2.0 [2022-04-27]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#29](https://github.com/turbot/steampipe-plugin-net/pull/29))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#27](https://github.com/turbot/steampipe-plugin-net/pull/27))

## v0.1.1 [2022-01-12]

_Enhancements_

- Updated the Slack channel links in the `docs/index.md` and the `README.md` files ([#18](https://github.com/turbot/steampipe-plugin-net/pull/18))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v183--2021-12-23) ([#19](https://github.com/turbot/steampipe-plugin-net/pull/19))

## v0.1.0 [2021-12-21]

_Enhancements_

- The `protocol` column of `net_connection` is now optional and it defaults to `tcp`

## v0.0.2 [2021-11-23]

_What's new?_

_Enhancements_

- Recompiled plugin with go version 1.17 ([#13](https://github.com/turbot/steampipe-plugin-net/pull/13))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) ([#12](https://github.com/turbot/steampipe-plugin-net/pull/12))

_Bug fixes_

- Fixed: SQL in first example of net_connection docs

## v0.0.1 [2021-04-28]

_What's new?_

- New tables added
  - [net_certificate](https://hub.steampipe.io/plugins/turbot/net/tables/net_certificate)
  - [net_connection](https://hub.steampipe.io/plugins/turbot/net/tables/net_connection)
  - [net_dns_record](https://hub.steampipe.io/plugins/turbot/net/tables/net_dns_record)
  - [net_dns_reverse](https://hub.steampipe.io/plugins/turbot/net/tables/net_dns_reverse)
