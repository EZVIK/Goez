package upload

import (
	"Goez/pkg/config"
	"Goez/pkg/files"
	"Goez/pkg/utils"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
)

func GetImageFullUrl(name string) string {
	return config.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

func GetImageName(name string) string {

	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = utils.Md5Crypt(fileName, utils.Upper(4))

	return fileName + ext
}

func GetImagePath() string {
	return config.AppSetting.ImageSavePath
}

func GetImageFullPath() string {
	return config.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {
	ext := files.GetExt(fileName)
	for _, allowExt := range config.AppSetting.ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := files.GetSize(f)
	if err != nil {
		log.Println(err)
		//logging.Warn(err)
		return false
	}

	return size <= config.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = files.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := files.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
