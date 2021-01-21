package config

const (
	CommonResFolder = "0"
	HeadImgFolder   = "1"
	AppApkFolder    = "2"
	IdCardFolder    = "3"
)

var FileDirMap = map[string]string{
	CommonResFolder: "./webroot/resource/",
	HeadImgFolder:   "./webroot/headImg/",
	AppApkFolder:    "./webroot/apk/",
	IdCardFolder:    "./webroot/idCard/",
}

func GetResourceDir(fileType string) string {
	path,ok := FileDirMap[fileType]
	if !ok{
		return FileDirMap[CommonResFolder]
	}
	return path
}
