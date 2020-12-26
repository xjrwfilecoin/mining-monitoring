package utils

import (
	"fmt"
	qrcode "github.com/skip2/go-qrcode"
	"os"
	"mining-monitoring/config"
	"time"
)

/**
生成二维码路径
*/
func GenerateQrCodePic(info string, username string) (string, error) {
	fileDir := config.QRCodePath
	ok, err := FileExists(fileDir)
	if err != nil {
		return "", err
	}
	if !ok {
		err := os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	fileName := fileDir + username + "_" + fmt.Sprintf("%d", time.Now().Nanosecond()) + ".png"

	err = qrcode.WriteFile(info, qrcode.Medium, 256, fileName)
	if err != nil {
		return "", err
	}
	return fileName, nil
}
