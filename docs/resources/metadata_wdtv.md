---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_metadata_wdtv Resource - terraform-provider-lidarr"
subcategory: "Metadata"
description: |-
  <!-- subcategory:Metadata -->
  
  Metadata Wdtv resource.
  For more information refer to Metadata https://wiki.servarr.com/lidarr/settings#metadata and WDTV https://wiki.servarr.com/lidarr/supported#wdtvmetadata.
---

# lidarr_metadata_wdtv (Resource)

<!-- subcategory:Metadata -->
Metadata Wdtv resource.
For more information refer to [Metadata](https://wiki.servarr.com/lidarr/settings#metadata) and [WDTV](https://wiki.servarr.com/lidarr/supported#wdtvmetadata).

## Example Usage

```terraform
resource "lidarr_metadata_wdtv" "example" {
  enable         = true
  name           = "Example"
  track_metadata = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Metadata name.
- `track_metadata` (Boolean) Track metadata flag.

### Optional

- `enable` (Boolean) Enable flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Metadata ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_metadata_wdtv.example 1
```
