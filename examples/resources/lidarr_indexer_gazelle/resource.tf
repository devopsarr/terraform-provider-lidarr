resource "lidarr_indexer_gazelle" "example" {
  enable_automatic_search = true
  name                    = "Example"
  base_url                = "https://orpheus.network"
  username                = "User"
  password                = "Pass"
  use_freelech_token      = false
  minimum_seeders         = 1
}
