package dtos

type RedisConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type RazorPayConfig struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

type WebAppConfigFile struct {
	Dynamo   map[string]string      `json:"dynamo"`
	Impulse  map[string]string      `json:"impulse"`
	WebApp   map[string]string      `json:"webapp"`
	MongoURL map[string]string      `json:"mongoURL"`
	Redis    map[string]RedisConfig `json:"redis"`
	SMTP     struct {
		Host string `json:"host"`
		Port string `json:"port"`
		Auth struct {
			User string `json:"user"`
			Pass string `json:"pass"`
		} `json:"auth"`
	} `json:"smtp"`
	RazorPay map[string]RazorPayConfig `json:"razorpay"`
	RootUrl  map[string]string         `json:"rootUrl"`
	Ingress  map[string]Ingress        `json:"Ingress"`
	Paymeter Paymeter                  `json:"paymeter"`
}
