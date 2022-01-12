package atlassian

type AtlassianConfig struct {
	BaseUrl    string `yaml:"base_url"`
	Username   string `yaml:"username"`
	Token      string `yaml:"token"`
	ProjectKey string `yaml:"project"`
}
