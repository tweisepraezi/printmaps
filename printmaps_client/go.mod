module github.com/printmaps/printmaps/printmaps_client

go 1.26

require (
	github.com/StefanSchroeder/Golang-Ellipsoid v0.0.0-20260222131200-12cd88ff1913
	github.com/davecgh/go-spew v1.1.1
	github.com/im7mortal/UTM v1.4.0
	github.com/paulmach/orb v0.12.0
	github.com/printmaps/printmaps/pd v1.0.0
	github.com/yuin/gopher-lua v1.1.1
	gopkg.in/yaml.v2 v2.4.0
)

require go.mongodb.org/mongo-driver v1.17.9 // indirect

replace github.com/printmaps/printmaps/pd => ../pd
