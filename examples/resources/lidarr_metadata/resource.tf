resource "lidarr_metadata" "example" {
  enable          = true
  name            = "Example"
  implementation  = "WdtvMetadata"
  config_contract = "WdtvMetadataSettings"
  track_metadata  = true
  tags            = [1, 2]
}