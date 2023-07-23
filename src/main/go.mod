module main

go 1.19

replace args => ../args

replace securitymodel => ../securitymodel

replace addb => ../addb

replace libadm => ../../../adm/src/libadm

require args v0.0.0-00010101000000-000000000000

require (
	github.com/cucumber/gherkin-go/v19 v19.0.3 // indirect
	github.com/cucumber/messages-go/v16 v16.0.1 // indirect
	github.com/fogleman/gg v1.3.0 // indirect
	github.com/goccy/go-graphviz v0.1.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/image v0.5.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	addb v0.0.0-00010101000000-000000000000 // indirect
	libadm v0.0.0-00010101000000-000000000000 // indirect
	securitymodel v0.0.0-00010101000000-000000000000 // indirect
)
