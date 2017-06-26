---
layout: "papertrail"
page\_title: "Papertrail: papertrail\_group"
sidebar\_current: "docs-papertrail-resource-papertrail\_group"
description: |-
  Manages a Papertrail group.
---

# papertrail_group

The ``papertrail_group`` resource creates and manages a
[Papertrail Group](http://help.papertrailapp.com/kb/how-it-works/settings-api/#groups).

## Usage

```hcl
resource "papertrail_group" "usage" {
  name = "Group Name"
  system_wildcard = "*/some_wild_card/*"
}
```

## Argument Reference

* `name` - (Required) Name of the group.

* `system_wildcard` - (Optional) wildcard for system names that belong to the group.

## Out Parameters

* `id` - Group ID
