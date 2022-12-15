resource "lidarr_indexer_waffles" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://www.waffles.ch"
  user_id                 = "User"
  rss_passkey             = "Pass"
  use_freelech_token      = false
  minimum_seeders         = 1
}
