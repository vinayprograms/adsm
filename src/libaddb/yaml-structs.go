package libaddb

type ItemType string
const (
	Human ItemType = "human"
	Program	ItemType = "program"
	System	ItemType = "system"
	Flow	ItemType = "flow"
)

type ADDBComponent struct {
	Id string `yaml:"id"`
	Name string `yaml:"name"`
	Type ItemType `yaml:"type"`
	Description string `yaml:"description"`
	DesignDocument string `yaml:"design-document"`
	Base []string `yaml:"base"`
	Recommendations []string `yaml:"recommendations"`
	ADM []string `yaml:"adm"`

	// Only for humans
	Roles []string `yaml:"roles"`
	Interface string `yaml:"interface"`

	// Only for programs/systems
	CodeRepository string `yaml:"repo"`
	Languages []string `yaml:"languages"`
	Dependencies []string	`yaml:"dependencies"`

	// Only for flows
	Protocol []string `yaml:"protocol"`
}