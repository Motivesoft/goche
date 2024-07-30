package uci

// These may be replaced at build time
var engineName string = "goche"
var authorName string = "Ian Brown"
var versionName string = "unknown"

func GetEngineName() string {
	return engineName
}

func GetAuthorName() string {
	return authorName
}

func GetVersionName() string {
	return versionName
}
