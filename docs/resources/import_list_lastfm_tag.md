---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_import_list_lastfm_tag Resource - terraform-provider-lidarr"
subcategory: "Import Lists"
description: |-
  <!-- subcategory:Import Lists -->
  
  Import List Last.fm Tag resource.
  For more information refer to Import List https://wiki.servarr.com/lidarr/settings#import-lists and Last.fm Tag https://wiki.servarr.com/lidarr/supported#lastfmtag.
---

# lidarr_import_list_lastfm_tag (Resource)

<!-- subcategory:Import Lists -->
Import List Last.fm Tag resource.
For more information refer to [Import List](https://wiki.servarr.com/lidarr/settings#import-lists) and [Last.fm Tag](https://wiki.servarr.com/lidarr/supported#lastfmtag).

## Example Usage

```terraform
resource "lidarr_import_list_lastfm_tag" "example" {
  enable_automatic_add = false
  should_monitor       = "specificAlbum"
  should_search        = false
  root_folder_path     = "/config"
  monitor_new_items    = "all"
  quality_profile_id   = 1
  metadata_profile_id  = 1
  name                 = "Example"
  tag_id               = "TagExample"
  count                = 25
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `count_list` (Number) Elements to pull from list.
- `name` (String) Import List name.
- `tag_id` (String) Tag ID.

### Optional

- `enable_automatic_add` (Boolean) Enable automatic add flag.
- `list_order` (Number) List order.
- `metadata_profile_id` (Number) Metadata profile ID.
- `monitor_new_items` (String) Monitor new items.
- `quality_profile_id` (Number) Quality profile ID.
- `root_folder_path` (String) Root folder path.
- `should_monitor` (String) Should monitor.
- `should_monitor_existing` (Boolean) Should monitor existing flag.
- `should_search` (Boolean) Should search flag.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Import List ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_import_list_lastfm_tag.example 1
```
