---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_notification_synology_indexer Resource - terraform-provider-lidarr"
subcategory: "Notifications"
description: |-
  Notification Synology Indexer resource.
  For more information refer to Notification https://wiki.servarr.com/lidarr/settings#connect and Synology https://wiki.servarr.com/lidarr/supported#synologyindexer.
---

# lidarr_notification_synology_indexer (Resource)

<!-- subcategory:Notifications -->Notification Synology Indexer resource.
For more information refer to [Notification](https://wiki.servarr.com/lidarr/settings#connect) and [Synology](https://wiki.servarr.com/lidarr/supported#synologyindexer).

## Example Usage

```terraform
resource "lidarr_notification_synology_indexer" "example" {
  on_upgrade        = true
  on_rename         = false
  on_track_retag    = false
  on_release_import = true

  name = "Example"

  update_library = true
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) NotificationSynology name.
- `on_release_import` (Boolean) On movie file delete for upgrade flag.
- `on_rename` (Boolean) On rename flag.
- `on_track_retag` (Boolean) On movie file delete flag.
- `on_upgrade` (Boolean) On upgrade flag.

### Optional

- `tags` (Set of Number) List of associated tags.
- `update_library` (Boolean) Update library flag.

### Read-Only

- `id` (Number) Notification ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_notification_synology_indexer.example 1
```