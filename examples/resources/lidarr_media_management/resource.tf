resource "lidarr_media_management" "example" {
  unmonitor_previous_tracks = true
  hardlinks_copy            = true
  create_empty_folders      = true
  delete_empty_folders      = true
  watch_library_for_changes = true
  import_extra_files        = true
  set_permissions           = true
  skip_free_space_check     = true
  minimum_free_space        = 100
  recycle_bin_days          = 7
  chmod_folder              = "755"
  chown_group               = "arrs"
  download_propers_repacks  = "preferAndUpgrade"
  allow_fingerprinting      = "never"
  extra_file_extensions     = "info"
  file_date                 = "none"
  recycle_bin_path          = ""
  rescan_after_refresh      = "always"
}