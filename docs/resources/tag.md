---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_tag Resource - terraform-provider-lidarr"
subcategory: ""
description: |-
  Tag resource
---

# lidarr_tag (Resource)

Tag resource

## Example Usage

```terraform
resource "lidarr_tag" "example" {
  label = "some-value"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `label` (String) Tag value

### Read-Only

- `id` (Number) Tag ID

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_tag.example 10
```
