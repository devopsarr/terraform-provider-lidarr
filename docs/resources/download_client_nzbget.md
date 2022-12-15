---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_download_client_nzbget Resource - terraform-provider-lidarr"
subcategory: "Download Clients"
description: |-
  Download Client NZBGet resource.
  For more information refer to Download Client https://wiki.servarr.com/lidarr/settings#download-clients and NZBGet https://wiki.servarr.com/lidarr/supported#nzbget.
---

# lidarr_download_client_nzbget (Resource)

<!-- subcategory:Download Clients -->Download Client NZBGet resource.
For more information refer to [Download Client](https://wiki.servarr.com/lidarr/settings#download-clients) and [NZBGet](https://wiki.servarr.com/lidarr/supported#nzbget).

## Example Usage

```terraform
resource "lidarr_download_client_nzbget" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "nzbget"
  url_base = "/nzbget/"
  port     = 6789
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.

### Optional

- `add_paused` (Boolean) Add paused flag.
- `enable` (Boolean) Enable flag.
- `host` (String) host.
- `music_category` (String) TV category.
- `older_music_priority` (Number) Older TV priority. `-100` VeryLow, `-50` Low, `0` Normal, `50` High, `100` VeryHigh, `900` Force.
- `password` (String, Sensitive) Password.
- `port` (Number) Port.
- `priority` (Number) Priority.
- `recent_music_priority` (Number) Recent TV priority. `-100` VeryLow, `-50` Low, `0` Normal, `50` High, `100` VeryHigh, `900` Force.
- `remove_completed_downloads` (Boolean) Remove completed downloads flag.
- `remove_failed_downloads` (Boolean) Remove failed downloads flag.
- `tags` (Set of Number) List of associated tags.
- `url_base` (String) Base URL.
- `use_ssl` (Boolean) Use SSL flag.
- `username` (String) Username.

### Read-Only

- `id` (Number) Download Client ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_download_client_nzbget.example 1
```