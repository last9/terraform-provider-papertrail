---
layout: "papertrail"
page\_title: "Provider: Papertrail"
sidebar_current: "docs-papertrail-index"
description: |-
  A provider for Papertrail.
---

# Papertrail Provider

The Papertrail provider gives the ability to manage a Papertrail account settings.

Use the navigation to the left to read about the available resources.

## Usage

```hcl
provider "papertrail" {
  token = "l5cqHkg23bhyqAGpgJ6R"
}
```

## Argument Reference

The following arguments are supported:

* `token` - (Required) The Papertrail API token.
