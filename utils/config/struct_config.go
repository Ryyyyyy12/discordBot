package config

type Config struct {
	BotToken string `yaml:"bot_token"`
	AuthKey  string `yaml:"auth_key"`
	Text1    string `yaml:"text_1"`
	Text2    string `yaml:"text_2"`
}
