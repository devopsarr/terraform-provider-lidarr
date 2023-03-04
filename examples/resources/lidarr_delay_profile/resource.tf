resource "lidarr_delay_profile" "example" {
  enable_usenet      = true
  enable_torrent     = true
  usenet_delay       = 0
  torrent_delay      = 0
  tags               = [1, 2]
  preferred_protocol = "torrent"
}