package common

//BotConfig contains configuration for the bot
type BotConfig struct {
	BotKey string `json:"bot_key"`
}

//Questions questions list
type Questions struct {
	Question string `json:"question"`
	Options  []string
	Answer   int
}
