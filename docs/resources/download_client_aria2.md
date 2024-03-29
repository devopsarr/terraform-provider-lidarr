---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_download_client_aria2 Resource - terraform-provider-lidarr"
subcategory: "Download Clients"
description: |-
  <!-- subcategory:Download Clients -->
  
  Download Client Aria2 resource.
  For more information refer to Download Client https://wiki.servarr.com/lidarr/settings#download-clients and Aria2 https://wiki.servarr.com/lidarr/supported#aria2.
---

# lidarr_download_client_aria2 (Resource)

<!-- subcategory:Download Clients -->
Download Client Aria2 resource.
For more information refer to [Download Client](https://wiki.servarr.com/lidarr/settings#download-clients) and [Aria2](https://wiki.servarr.com/lidarr/supported#aria2).

## Example Usage

```terraform
resource "lidarr_download_client_aria2" "example" {
  enable   = true
  priority = 1
  name     = "Example"
  host     = "aria2"
  rpc_path = "/aria2/"
  port     = 6800
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Download Client name.

### Optional

- `enable` (Boolean) Enable flag.
- `host` (String) host.
- `port` (Number) Port.
- `priority` (Number) Priority.
- `remove_completed_downloads` (Boolean) Remove completed downloads flag.
- `remove_failed_downloads` (Boolean) Remove failed downloads flag.
- `rpc_path` (String) RPC path.
- `secret_token` (String) Secret token.
- `tags` (Set of Number) List of associated tags.
- `use_ssl` (Boolean) Use SSL flag.

### Read-Only

- `id` (Number) Download Client ID.

## Import

Import is supported using the following syntax:

```shell
# import using the API/UI ID
terraform import lidarr_download_client_aria2.example 1
```
