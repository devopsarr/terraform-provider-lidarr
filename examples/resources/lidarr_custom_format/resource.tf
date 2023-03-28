resource "lidarr_custom_format" "example" {
  include_custom_format_when_renaming = true
  name                                = "Example"

  specifications = [
    {
      name           = "Preferred Words"
      implementation = "ReleaseTitleSpecification"
      negate         = false
      required       = false
      value          = "\\b(SPARKS|Framestor)\\b"
    },
    {
      name           = "Size"
      implementation = "SizeSpecification"
      negate         = false
      required       = false
      min            = 0
      max            = 100
    }
  ]
}