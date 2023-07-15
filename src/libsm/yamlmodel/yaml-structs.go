package yamlmodel

type SecurityModel struct {
	Title string				`yaml:"title"`
	DesignDocument string `yaml:"design-document"`
	AddbUri string `yaml:"addb"`
	ModelADM []string `yaml:"adm,flow"`
	Externals []*Entity `yaml:"externals,flow"`
	Entities []*Entity `yaml:"entities,flow"`
	Flows []*Flow `yaml:"flows,flow"`

	// internal variable to locate adm
	AdmDir string
}

type ItemType string
const (
	Human ItemType = "human"
	Program	ItemType = "program"
	System	ItemType = "system"
	Role		ItemType = "role"
)

type Entity struct {
	Id string `yaml:"id"`
	Type ItemType `yaml:"type"`
	Name string `yaml:"name"`
	Description string `yaml:"description"`
	Base []string `yaml:"base"`
	Mitigations []string `yaml:"mitigations"`	// Not applicable for external entities
	Recommendations []string `yaml:"recommendations"`	// Not applicable for external entities
	ADM []string `yaml:"adm"`							// Not applicable for external entities
	
	// Only for humans
	Interface string `yaml:"interface"`

	// Only for programs/systems
	Roles []string `yaml:"roles"`
	CodeRepository string `yaml:"repo"`							// Not applicable for external entities
	Languages []string `yaml:"languages"`
	Dependencies []string	`yaml:"dependencies"`			// Not applicable for external entities

	// internal variable to locate adm
	AdmDir string
}

type Flow struct {
	Id string `yaml:"id"`
	Name string `yaml:"name"`
	Description string `yaml:"description"`
	Protocol []string `yaml:"protocol"`
	Sender string `yaml:"sender"`
	Receiver string `yaml:"receiver"`
	Mitigations []string `yaml:"mitigations"`
	Recommendations []string `yaml:"recommendations"`
	ADM []string `yaml:"adm"`

	// internal variable to locate adm
	AdmDir string
}