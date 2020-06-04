package env

import "os"

type Config struct {
	ProjectID      string
	DialogflowAuth string
	DatastoreAuth  string
	Language       string
	Timezone       string
}

func Get() *Config {
	return &Config{
		ProjectID:      os.Getenv("PROJECT_ID"),
		DialogflowAuth: os.Getenv("DIALOGFLOW_KEYFILE_JSON"),
		DatastoreAuth:  os.Getenv("DATASTORE_CRED_JSON"),
		Language:       os.Getenv("LANGUAGE"),
		Timezone:       os.Getenv("TIMEZONE"),
	}
}
