package config

// Для передачи сообщений через kafka в indexer
const (
	DaoIndexerTopic = "dao_indexer"
	DaoIndexerGroup = "dao_indexer-controller_group"
)

const (
	DaoControllerBotTopic = "dao_bot"
	DaoControllerBotGroup = "dao_controller-bot_group"
)
