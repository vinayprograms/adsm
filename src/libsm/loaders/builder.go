package loaders

type Builder struct {
	yamlIndex map[string]interface{}
	objectIndex map[string]interface{}
}

func (t *Builder) init() {
	if t.yamlIndex == nil {
		t.yamlIndex = make(map[string]interface{})
	}
	if t.objectIndex == nil {
		t.objectIndex = make(map[string]interface{})
	}
}