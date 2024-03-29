---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_release_profile Resource - terraform-provider-lidarr"
subcategory: "Profiles"
description: |-
  <!-- subcategory:Profiles -->
  
  Release Profile resource.
  For more information refer to Release Profiles https://wiki.servarr.com/lidarr/settings#release-profiles documentation.
---

# lidarr_release_profile (Resource)

<!-- subcategory:Profiles -->
Release Profile resource.
For more information refer to [Release Profiles](https://wiki.servarr.com/lidarr/settings#release-profiles) documentation.

## Example Usage

```terraform
resource "lidarr_release_profile" "example" {
  enabled                         = true
  include_preferred_when_renaming = true
  indexer_id                      = 0
  required                        = "dolby,digital"
  ignored                         = "mp3"
  preferred = [
    {
      term  = "higher"
      score = 100
    },
    {
      term  = "lower"
      score = -100
    },
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `enabled` (Boolean) Enabled.
- `ignored` (Set of String) Ignored terms. At least one of `required` and `ignored` must be set.
- `indexer_id` (Number) Indexer ID. Default to all.
- `required` (Set of String) Required terms. At least one of `required` and `ignored` must be set.
- `tags` (Set of Number) List of associated tags.

### Read-Only

- `id` (Number) Release Profile ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_release_profile.example 10
```
