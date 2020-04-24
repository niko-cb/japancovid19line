package env

import "os"

func ProjectID() string {
	return os.Getenv("PROJECT_ID")
}

func AuthDialogflow() string {
	return os.Getenv("DIALOGFLOW_KEYFILE_JSON")
}

func AuthDatastore() string {
	return os.Getenv("DATASTORE_CRED_JSON")
}

func Language() string {
	return os.Getenv("LANGUAGE")
}

func Timezone() string {
	return os.Getenv("TIMEZONE")
}
