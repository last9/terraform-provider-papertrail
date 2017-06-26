---
layout: "papertrail"
page\_title: "Papertrail: papertrail\_system\_group"
sidebar\_current: "docs-papertrail-resource-papertrail\_system\_group"
description: |-
  Manages a Papertrail systems in a group.
---

# papertrail\_system_group

The ``papertrail_system_group`` resource to manage group systems.
[Papertrail Group Systems](http://help.papertrailapp.com/kb/how-it-works/settings-api/#join-group).

## Usage

```hcl
resource "papertrail_system_group" "sg" {
  system_id = 321
  group_id = 123
}
```

## Argument Reference

* `system_id` - (Required) Papertrail System ID.

* `group_id` - (Required) Papertrail Group ID

## Out Parameters

* `index` - Index of the system in systems list of the group.
