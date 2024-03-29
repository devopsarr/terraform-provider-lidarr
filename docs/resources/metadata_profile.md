---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_metadata_profile Resource - terraform-provider-lidarr"
subcategory: "Profiles"
description: |-
  <!-- subcategory:Profiles -->
  
  Metadata Profile resource.
  For more information refer to Metadata Profile https://wiki.servarr.com/lidarr/settings#metadata-profiles documentation.
---

# lidarr_metadata_profile (Resource)

<!-- subcategory:Profiles -->
Metadata Profile resource.
For more information refer to [Metadata Profile](https://wiki.servarr.com/lidarr/settings#metadata-profiles) documentation.

## Example Usage

```terraform
resource "lidarr_metadata_profile" "example" {
  name                  = "Example"
  primary_album_types   = [1, 2]
  secondary_album_types = [1]
  release_statuses      = [3]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Metadata Profile name.
- `primary_album_types` (Set of Number) Primary album types.
- `release_statuses` (Set of Number) Release statuses.
- `secondary_album_types` (Set of Number) Secondary album types.

### Read-Only

- `id` (Number) Metadata Profile ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_metadata_profile.example 10
```
