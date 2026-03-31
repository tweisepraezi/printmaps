module github.com/printmaps/printmaps/printmaps_buildservice

go 1.26

require (
	github.com/printmaps/printmaps/pd v1.0.0
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/printmaps/printmaps/pd => ../pd
