---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_custom_format Data Source - terraform-provider-lidarr"
subcategory: "Profiles"
description: |-
  <!-- subcategory:Profiles -->
  
  Single Custom Format ../resources/custom_format.
---

# lidarr_custom_format (Data Source)

<!-- subcategory:Profiles -->
Single [Custom Format](../resources/custom_format).

## Example Usage

```terraform
data "lidarr_custom_format" "example" {
  name = "Example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Custom Format name.

### Read-Only

- `id` (Number) Custom Format ID.
- `include_custom_format_when_renaming` (Boolean) Include custom format when renaming flag.
- `specifications` (Attributes Set) Specifications. (see [below for nested schema](#nestedatt--specifications))

<a id="nestedatt--specifications"></a>
### Nested Schema for `specifications`

Read-Only:

- `implementation` (String) Implementation.
- `max` (Number) Max.
- `min` (Number) Min.
- `name` (String) Specification name.
- `negate` (Boolean) Negate flag.
- `required` (Boolean) Computed flag.
- `value` (String) Value.


