package types

type ConfigurationPlugin struct {
	Debug    bool   `envconfig:"PLUGIN_DEBUG"       default:"false"`
	ProxyURL string `envconfig:"PLUGIN_PROXY_URL"   default:"https://api.telegram.org"`
	Token    string `envconfig:"PLUGIN_TOKEN"`
	ChatID   int    `envconfig:"PLUGIN_CHAT_ID"`
}

type ConfigurationDrone struct {
	BuildStatus   string `envconfig:"DRONE_BUILD_STATUS"`
	BuildNumber   string `envconfig:"DRONE_BUILD_NUMBER"`
	BuildLink     string `envconfig:"DRONE_BUILD_LINK"`
	RepoName      string `envconfig:"DRONE_REPO_NAME"`
	CommitMessage string `envconfig:"DRONE_COMMIT_MESSAGE"`
}

type MessageBody struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}
