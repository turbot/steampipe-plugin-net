
# Steampipe Net Plugin

The `net` plugin covers various network requests and utilities that are helpful
for checking and querying network components.

## Examples

```
select * from net_certificate where domain_name = 'steampipe.io';
```

### Joining with a domains list

TODO - Certificate list needs to accept the quals so it can decide which rows to support

```
select name, (select subject from net.net_certificate where domain_name = name) from domains;
```

```
select d.*, c.* from (select * from domains) d left join lateral (select * from net.net_certificate where domain_name = d.name) c on true;
```

```
select * from domains as d left join net.net_certificate as c on d.name = c.domain_name;
```
