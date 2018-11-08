package dtos

type RedisConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type WebAppConfigFile struct {
	Dynamo map[string]string `json:"dynamo"`
	Impulse map[string]string `json:"impulse"`
	WebApp map[string]string `json:"webapp"`
	MongoURL map[string]string 	`json:"mongoURL"`
	Redis map[string]RedisConfig `json:"redis"`
}