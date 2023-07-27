![image](https://hub.steampipe.io/images/plugins/turbot/net-social-graphic.png)

# Net Plugin for Steampipe

Use SQL to query DNS records, certificates and other network information. Open source CLI. No DB required.

* **[Get started →](https://hub.steampipe.io/plugins/turbot/net)**
* Documentation: [Table definitions & examples](https://hub.steampipe.io/plugins/turbot/net/tables)
* Community: [Join #steampipe on Slack →](https://turbot.com/community/join)
* Get involved: [Issues](https://github.com/turbot/steampipe-plugin-net/issues)

## Quick start

Install the plugin with [Steampipe](https://steampipe.io):
```shell
steampipe plugin install net
```

Run a query:
```sql
select * from net_certificate where domain = 'steampipe.io';
```

## Developing

Prerequisites:
- [Steampipe](https://steampipe.io/downloads)
- [Golang](https://golang.org/doc/install)

Clone:

```sh
git clone https://github.com/turbot/steampipe-plugin-net.git
cd steampipe-plugin-net
```

Build, which automatically installs the new version to your `~/.steampipe/plugins` directory:
```
make
```

Configure the plugin:
```
cp config/* ~/.steampipe/config
```

Try it!
```
steampipe query
> .inspect net
```

Further reading:
* [Writing plugins](https://steampipe.io/docs/develop/writing-plugins)
* [Writing your first table](https://steampipe.io/docs/develop/writing-your-first-table)

## Contributing

Please see the [contribution guidelines](https://github.com/turbot/steampipe/blob/main/CONTRIBUTING.md) and our [code of conduct](https://github.com/turbot/steampipe/blob/main/CODE_OF_CONDUCT.md). All contributions are subject to the [Apache 2.0 open source license](https://github.com/turbot/steampipe-plugin-net/blob/main/LICENSE).

`help wanted` issues:
- [Steampipe](https://github.com/turbot/steampipe/labels/help%20wanted)
- [Net Plugin](https://github.com/turbot/steampipe-plugin-net/labels/help%20wanted)
