## v0.10.0 [2023-09-12]

_Deprecations_

- Deprecated `domain` column in `net_certificate` table, which has been replaced by the `address` column. Please note that the `address` column requires a port, e.g., `github.com:443`. This column will be removed in a future version. ([#50](https://github.com/turbot/steampipe-plugin-net/pull/50))

_What's new?_

- Added `address` column to the `net_certificate` table to allow specifying a port with the domain name. ([#50](https://github.com/turbot/steampipe-plugin-net/pull/50))

## v0.9.0 [2023-04-10]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#58](https://github.com/turbot/steampipe-plugin-net/pull/58))

## v0.8.1 [2022-12-27]

_Bug fixes_

- Fixed the example query in `docs/index.md` to correctly check the expiry date of certificates associated with websites. ([#53](https://github.com/turbot/steampipe-plugin-net/pull/53)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Fixed the typo in the example query of `net_connection` table document. ([#54](https://github.com/turbot/steampipe-plugin-net/pull/54)) (Thanks [@muescha](https://github.com/muescha) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.8](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v418-2022-09-08) which increases the default open file limit. ([#52](https://github.com/turbot/steampipe-plugin-net/pull/52))

## v0.8.0 [2022-09-28]

_Bug fixes_

- Fixed the `net_certificate` table to return `null` instead of an error when the queried host doesn't exist in a DNS. ([#49](https://github.com/turbot/steampipe-plugin-net/pull/49)) (Thanks [@bdd4329](https://github.com/bdd4329) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v417-2022-09-08) which includes several caching and memory management improvements. ([#46](https://github.com/turbot/steampipe-plugin-net/pull/46))
- Recompiled plugin with Go version `1.19`. ([#46](https://github.com/turbot/steampipe-plugin-net/pull/46))

## v0.7.0 [2022-08-17]

_Bug fixes_

- Fixed the `net_certificate` table to return an empty row instead of an error if the domain does not have an SSL certificate. ([#42](https://github.com/turbot/steampipe-plugin-net/pull/42))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v332--2022-07-11) which includes several caching fixes. ([#39](https://github.com/turbot/steampipe-plugin-net/pull/39))

## v0.6.0 [2022-06-30]

_Enhancements_

- Improved the backoff/retry mechanism in the `net_dns_record` table to return results faster and more reliably. ([#32](https://github.com/turbot/steampipe-plugin-net/pull/32))
- Recompiled plugin with [miekg/dns v1.1.50](https://github.com/miekg/dns/releases/tag/v1.1.50). ([#32](https://github.com/turbot/steampipe-plugin-net/pull/32))

## v0.5.0 [2022-06-24]

_What's new?_

- New tables added:
  - [net_http_request](https://hub.steampipe.io/plugins/turbot/net/tables/net_http_request) ([#38](https://github.com/turbot/steampipe-plugin-net/pull/38))

_Enhancements_

- Updated the `net_certificate` table to return an error for any non-200 response status codes. ([#37](https://github.com/turbot/steampipe-plugin-net/pull/37))
- Recompiled plugin with [steampipe-plugin-sdk v3.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v330--2022-6-22). ([#38](https://github.com/turbot/steampipe-plugin-net/pull/38))

## v0.4.0 [2022-06-09]

_What's new?_

- New tables added:
  - [net_tls_connection](https://hub.steampipe.io/plugins/turbot/net/tables/net_tls_connection) ([#35](https://github.com/turbot/steampipe-plugin-net/pull/35))

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
