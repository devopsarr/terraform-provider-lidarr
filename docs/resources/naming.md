---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_naming Resource - terraform-provider-lidarr"
subcategory: "Media Management"
description: |-
  <!-- subcategory:Media Management -->
  
  Naming resource.
  For more information refer to Naming https://wiki.servarr.com/lidarr/settings#community-naming-suggestions documentation.
---

# lidarr_naming (Resource)

<!-- subcategory:Media Management -->
Naming resource.
For more information refer to [Naming](https://wiki.servarr.com/lidarr/settings#community-naming-suggestions) documentation.

## Example Usage

```terraform
resource "lidarr_naming" "example" {
  rename_tracks              = true
  replace_illegal_characters = true
  standard_track_format      = "{Album Title} ({Release Year})/{Artist Name} - {Album Title} - {track:00} - {Track Title}"
  multi_disc_track_format    = "{Album Title} ({Release Year})/{Medium Format} {medium:00}/{Artist Name} - {Album Title} - {track:00} - {Track Title}"
  artist_folder_format       = "{Artist Name}"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `artist_folder_format` (String) Artist folder format.
- `multi_disc_track_format` (String) Multi disc track format.
- `rename_tracks` (Boolean) Lidarr will use the existing file name if false.
- `replace_illegal_characters` (Boolean) Replace illegal characters. They will be removed if false.
- `standard_track_format` (String) Standard track formatss.

### Read-Only

- `id` (Number) Naming ID.

## Import

Import is supported using the following syntax:

```shell
# import
terraform import lidarr_naming.example
```
