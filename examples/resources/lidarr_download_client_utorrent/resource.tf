resource "lidarr_download_client_utorrent" "example" {
  enable         = true
  priority       = 1
  name           = "Example"
  host           = "utorrent"
  url_base       = "/utorrent/"
  port           = 9091
  music_category = "tv-lidarr"
}