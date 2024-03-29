---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_metadata_kodi Resource - terraform-provider-lidarr"
subcategory: "Metadata"
description: |-
  <!-- subcategory:Metadata -->
  
  Metadata Kodi resource.
  For more information refer to Metadata https://wiki.servarr.com/lidarr/settings#metadata and KODI https://wiki.servarr.com/lidarr/supported#xbmcmetadata.
---

# lidarr_metadata_kodi (Resource)

<!-- subcategory:Metadata -->
Metadata Kodi resource.
For more information refer to [Metadata](https://wiki.servarr.com/lidarr/settings#metadata) and [KODI](https://wiki.servarr.com/lidarr/supported#xbmcmetadata).

## Example Usage

```terraform
resource "lidarr_metadata_kodi" "example" {
  enable          = true
  name            = "Example"
  artist_metadata = true
  album_images    = true
  artist_images   = true
  album_metadata  = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `album_images` (Boolean) Album images flag.
- `album_metadata` (Boolean) Album metadata flag.
- `artist_images` (Boolean) Artist images flag.
- `artist_metadata` (Boolean) Artist metadata flag.
- `name` (String) Metadata name.

### Optional

- `enable` (Boolean) Enable flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Metadata ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_metadata_kodi.example 1
```
