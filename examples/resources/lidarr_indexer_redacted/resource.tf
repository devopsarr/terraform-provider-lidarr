resource "lidarr_indexer_redacted" "example" {
  enable_automatic_search = true
  name                    = "Example"
  api_key                 = "Key"
  use_freelech_token      = false
  minimum_seeders         = 1
}
