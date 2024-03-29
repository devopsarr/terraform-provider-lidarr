---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_indexer_filelist Resource - terraform-provider-lidarr"
subcategory: "Indexers"
description: |-
  <!-- subcategory:Indexers -->
  
  Indexer FileList resource.
  For more information refer to Indexer https://wiki.servarr.com/lidarr/settings#indexers and FileList https://wiki.servarr.com/lidarr/supported#filelist.
---

# lidarr_indexer_filelist (Resource)

<!-- subcategory:Indexers -->
Indexer FileList resource.
For more information refer to [Indexer](https://wiki.servarr.com/lidarr/settings#indexers) and [FileList](https://wiki.servarr.com/lidarr/supported#filelist).

## Example Usage

```terraform
resource "lidarr_indexer_filelist" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://filelist.io"
  username                = "User"
  passkey                 = "PassKey"
  minimum_seeders         = 1
  categories              = [4, 6, 1]
  required_flags          = [1, 4]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) IndexerFilelist name.
- `passkey` (String, Sensitive) Passkey.
- `username` (String) Username.

### Optional

- `base_url` (String) Base URL.
- `categories` (Set of Number) Categories list.
- `enable_automatic_search` (Boolean) Enable automatic search flag.
- `enable_interactive_search` (Boolean) Enable interactive search flag.
- `enable_rss` (Boolean) Enable RSS flag.
- `minimum_seeders` (Number) Minimum seeders.
- `priority` (Number) Priority.
- `seed_ratio` (Number) Seed ratio.
- `seed_time` (Number) Seed time.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) IndexerFilelist ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_indexer_filelist.example 1
```
