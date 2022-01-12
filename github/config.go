package github

type GithubConfig struct {
	Token      string `yaml:"token"`
	Owner      string `yaml:"owner"`
	Repository string `yaml:"repo"`
}
