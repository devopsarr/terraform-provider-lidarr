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