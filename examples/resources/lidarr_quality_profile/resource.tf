resource "lidarr_quality_profile" "example" {
  name            = "example-lossless"
  upgrade_allowed = true
  cutoff          = 1100

  quality_groups = [
    {
      id   = 1100
      name = "lossless"
      qualities = [
        {
          id   = 7
          name = "ALAC"
        },
        {
          id   = 6
          name = "FLAC"
        }
      ]
    }
  ]
}