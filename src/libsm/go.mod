module libsm

go 1.19

replace libaddb => ../libaddb

replace libadm => ../../../adm/src/libadm

replace libsm/addb => ./addb

replace libsm/loaders => ./loaders

replace libsm/yamlmodel => ./yamlmodel

replace libsm/objmodel => ./objmodel

require gopkg.in/yaml.v3 v3.0.1 // direct

require libaddb v0.0.0-00010101000000-000000000000

require (
	github.com/cucumber/gherkin-go/v19 v19.0.3 // indirect
	github.com/cucumber/messages-go/v16 v16.0.1 // indirect
	github.com/gofrs/uuid v4.3.0+incompatible // indirect
	libadm v0.0.0-00010101000000-000000000000
)
