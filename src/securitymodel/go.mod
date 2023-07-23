module securitymodel

go 1.19

replace addb => ../addb

replace libadm => ../../../adm/src/libadm

replace securitymodel/addb => ./addb

replace securitymodel/loaders => ./loaders

replace securitymodel/yamlmodel => ./yamlmodel

replace securitymodel/objmodel => ./objmodel

require gopkg.in/yaml.v3 v3.0.1 // direct

require addb v0.0.0-00010101000000-000000000000

require (
	github.com/cucumber/gherkin-go/v19 v19.0.3 // indirect
	github.com/cucumber/messages-go/v16 v16.0.1 // indirect
	github.com/gofrs/uuid v4.3.0+incompatible // indirect
	libadm v0.0.0-00010101000000-000000000000
)
