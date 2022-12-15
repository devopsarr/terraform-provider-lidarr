resource "lidarr_indexer_headphones" "example" {
  enable_automatic_search = true
  name                    = "Example"
  username                = "User"
  password                = "Pass"
  categories              = [3000, 3010, 3020, 3030, 3040]
  tags                    = [1, 2]
}