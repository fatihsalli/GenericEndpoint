package configs

type Config struct {
	Server struct {
		Port map[string]string
		Host string
	}
	Database struct {
		Connection          string
		DatabaseName        string
		UserCollectionName  string
		OrderCollectionName string
	}
	Elasticsearch struct {
		Addresses map[string]string
		IndexName map[string]string
	}
}

var Configs = map[string]Config{
	"test": {
		Server: struct {
			Port map[string]string
			Host string
		}{
			Port: map[string]string{
				"orderAPI": ":8011",
			},
			Host: "localhost",
		},
		Database: struct {
			Connection          string
			DatabaseName        string
			UserCollectionName  string
			OrderCollectionName string
		}{
			Connection:          "mongodb://localhost:27017",
			DatabaseName:        "ProjectDB",
			UserCollectionName:  "Users",
			OrderCollectionName: "Orders",
		},
		Elasticsearch: struct {
			Addresses map[string]string
			IndexName map[string]string
		}{
			Addresses: map[string]string{
				"Address 1": "http://localhost:9200",
			},
			IndexName: map[string]string{
				"Order": "order_duplicate_v01",
			},
		},
	},
	"qa":   {},
	"prod": {},
}

func GetConfig(env string) Config {
	if conf, ok := Configs[env]; ok {
		return conf
	}

	return Configs["test"]
}
