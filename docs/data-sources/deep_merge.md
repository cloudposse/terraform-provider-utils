---
page_title: "deep_merge Data Source - terraform-provider-utils"
subcategory: ""
description: |-
  The deep_merge data source accepts a list of maps as input and deep merges them as output.
---

# Data Source `deep_merge`

The `deep_merge` data source accepts a list of maps as input and deep merges them as output.

## Example Usage

```terraform
data "deep_merge" "example" {
  inputs = {
    foo = "bar"
    baz = "bat"
  }
}
```

## Schema

### Required

- **inputs** (Map of String) A listx of arbitrary maps that is deep merged into the `output` attribute.

### Optional

- **id** (String) The ID of this resource.

### Read-only

- **output** (Map of String) The deep-merged map.


