package model

type RuntimeConfig struct {
	Debug         bool // 是否是debug
	LogPath       string
	HTTPListen    string
	MongodbUrl    string
	DevMongodbUrl string
}
