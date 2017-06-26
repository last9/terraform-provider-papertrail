---
layout: "papertrail"
page\_title: "Papertrail: papertrail\_system"
sidebar_current: "docs-circonus-resource-papertrail\_system"
description: |-
  Manages a papertrail System.
---

# papertrail_system

The ``papertrail_system`` resource creates and manages a
[Papertrail System](http://help.papertrailapp.com/kb/how-it-works/settings-api/#systems).


## Usage

```hcl
resource "papertrail_system" "my_system" {
  name = "My System"
  destination_id = 123
  destination_port = 8888
  ip_address = "127.0.0.1"
  hostname = "my.hostname.com"
  description = "Some Description"
}
```

## Argument Reference

* `name` - (Required) Name for papertrail system.

* `destination_id` = (Optional) Log destination id.

* `destination_port` = (Optional) Log destination port (One of destination `id` or `port` is Required).

* `ip_address` - (Optional) Ip address of the system.

* `hostname` - (Optional) Hostname of the system.

* `description` - (Optional) Description

## Out Parameters

The following attributes are exported:

* `id` - The System ID.

* `last_event_at` - timestamp at which last event happend.

* `syslog_hostname` - Log destination hostname.

* `syslog_port` - Log destination port.
