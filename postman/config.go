package postman

type Config struct {
	Name    string `yaml:"name"`
	Logging struct {
		Level       string `yaml:"level"`
		JSONLogging bool   `yaml:"jsonLogging"`
	} `yaml:"logging"`
	Port    int    `yaml:"port"`
	Mode    string `yaml:"mode"`
	Postman struct {
		token    string   `yaml:"apiToken"`
		URL         string   `yaml:"url"`
		Collections []string `yaml:collections`
	} `yaml:"postman"`
	Static struct {
		WatchFile bool `yaml:"watchFile"`
		Path string `yaml:"path"`
	} `yaml:"static"`
}

