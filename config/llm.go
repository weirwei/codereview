package config

type LLMConf struct {
	BaseUrl  string `json:"base_url"`
	SK       string `json:"sk"`
	Model    string `json:"model"`
	MaxToken string `json:"max_token"`
}
