---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "lidarr_metadata_consumers Data Source - terraform-provider-lidarr"
subcategory: "Metadata"
description: |-
  <!-- subcategory:Metadata -->
  
  List all available Metadata Consumers ../resources/metadata.
---

# lidarr_metadata_consumers (Data Source)

<!-- subcategory:Metadata -->
List all available [Metadata Consumers](../resources/metadata).

## Example Usage

```terraform
data "lidarr_metadata_consumers" "example" {
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `id` (String) The ID of this resource.
- `metadata_consumers` (Attributes Set) MetadataConsumer list. (see [below for nested schema](#nestedatt--metadata_consumers))

<a id="nestedatt--metadata_consumers"></a>
### Nested Schema for `metadata_consumers`

Read-Only:

- `album_images` (Boolean) Album images flag.
- `album_metadata` (Boolean) Album metadata flag.
- `artist_images` (Boolean) Artist images flag.
- `artist_metadata` (Boolean) Artist metadata flag.
- `config_contract` (String) Metadata configuration template.
- `enable` (Boolean) Enable flag.
- `id` (Number) Metadata ID.
- `implementation` (String) Metadata implementation name.
- `name` (String) Metadata name.
- `tags` (Set of Number) List of associated tags.
- `track_metadata` (Boolean) Track metadata flag.


