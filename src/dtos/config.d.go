package dtos

type HyperionConfig struct {
	IpfsPort        string `json:"ipfsPort"`
	IpfsClusterPort string `json:"ipfsClusterPort"`
}

type Auth struct {
	User  string `json:"user"`
	Pass  string `json:"pass"`
	Token string `json:"token"`
}

type Ingress struct {
	Annotations map[string]string `json:"Annotations"`
	Secret      string            `json:"secretName"`
}

type LocationConfig struct {
	MasterAPIHost    string         `json:"masterAPIHost"`
	WorkerNodeIP     string         `json:"workerNodeIP"`
	LocationCode     string         `json:"locationCode"`
	LocationName     string         `json:"locationName"`
	DynamoDomainName string         `json:"dynamoDomainName"`
	APIHost          string         `json:"apiHost"`
	Auth             Auth           `json:"auth"`
	Hyperion         HyperionConfig `json:"hyperion"`
}

type ClusterConfig struct {
	Clusters map[string]map[string]*LocationConfig `json:"clusters"`
}

type Blockchain struct {
	Testnet struct {
		URL string `json:"url"`
	} `json:"testnet"`
	Mainnet struct {
		URL string `json:"url"`
	} `json:"mainnet"`
}

type Paymeter struct {
	Blockchain map[string]Blockchain `json:"blockchains"`
	ApiKeys    map[string]string `json:"api_keys"`
}

type WebAppConfig struct {
	MongoConnectionURL string         `json:"mongoComongoURLnnectionURL"`
	RedisHost          string         `json:"redisHost"`
	RedisPort          string         `json:"redisPort"`
	ImageRepository    string         `json:"webAppImageName"`
	RazorPay           RazorPayConfig `json:"razorpay"`
	RootURL            string         `json:"rootUrl"`
	Ingress            Ingress        `json:"Ingress"`
	Paymeter           Paymeter       `json:"paymeter"`
}
