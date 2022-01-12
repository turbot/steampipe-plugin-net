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
