module github.com/printmaps/printmaps/printmaps_webservice

go 1.26

require (
	github.com/JamesMilnerUK/pip-go v0.0.0-20180711171552-99c4cbbc7deb
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/julienschmidt/httprouter v1.3.0
	github.com/printmaps/printmaps/pd v1.0.0
	github.com/rs/cors v1.11.1
	gopkg.in/yaml.v2 v2.4.0
)

require github.com/gabriel-vasile/mimetype v1.4.13

replace github.com/printmaps/printmaps/pd => ../pd
