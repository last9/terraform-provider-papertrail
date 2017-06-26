---
layout: "papertrail"
page\_title: "Papertrail: user"
sidebar\_current: "docs-papertrail-datasource-papertrail\_user"
description: |-
    Provides details about a specific Papertrail User.
---

# papertrail_user

`papertrail_user` provides
[details](http://help.papertrailapp.com/kb/how-it-works/settings-api/#users) about a specific
[User](http://help.papertrailapp.com/kb/how-it-works/settings-api/#users).

The `papertrail_user` data source can be used for pulling various attributes
about a specific Papertrail User. As well as validating a given email id.

## Example Usage

The following example shows how the resource might be used to obtain the metrics
usage and limit of a given Papertrail User.

```hcl
data "papertrail_user" "my_user" {
  email = "abc@xyz.com"
}
```

## Argument Reference

The arguments of this data source act as filters for querying the available
regions. The given filters must match exactly one region whose data will be
exported as attributes.

* `email` - (Required) The email id of the user.

## Attributes Reference

The following attributes are exported:

* `id` - Papertrail user id.

* `email` - Email id.
