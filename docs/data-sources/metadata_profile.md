---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_metadata_profile Data Source - terraform-provider-lidarr"
subcategory: "Profiles"
description: |-
  <!-- subcategory:Profiles -->
  
  Single Metadata Profile ../resources/metadata_profile.
---

# lidarr_metadata_profile (Data Source)

<!-- subcategory:Profiles -->
Single [Metadata Profile](../resources/metadata_profile).

## Example Usage

```terraform
data "lidarr_metadata_profile" "example" {
  name = "Example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Metadata Profile name.

### Read-Only

- `id` (Number) Metadata Profile ID.
- `primary_album_types` (Set of Number) Primary album types.
- `release_statuses` (Set of Number) Release statuses.
- `secondary_album_types` (Set of Number) Secondary album types.


