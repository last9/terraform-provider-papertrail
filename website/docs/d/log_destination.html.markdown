---
layout: "papertrail"
page\_title: "Papertrail: log\_destination"
sidebar\_current: "docs-circonus-datasource-papertrail\_log_destination"
description: |-
    Provides details about log destination.
---

# papertrail\_log_destination

`papertrail_log_destination` provides
[details](http://help.papertrailapp.com/kb/how-it-works/settings-api/#destinations) about a specific
[Log Destination](http://help.papertrailapp.com/kb/how-it-works/settings-api/#destinations).

## Example Usage

The following example shows how the data source might be used to obtain
the destination_id using destination port of any papertrail log destination.

```hcl
data "papertrail_log_destination" "ld" {
  port = 36267
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
regions. The given filters must match exactly one region whose data will be
exported as attributes.

* `id` - (Optional) The Destination ID of a given log destination.
* `port` - (Port) The Port of a given log destination.

At least one of the above attributes should be provided when searching for a destination.

## Attributes Reference

The following attributes are exported:

* `id` - The Destination ID.

* `port` - The Destination Port.

* `hostname` - The Hostname of destination.

* `description` - Description
