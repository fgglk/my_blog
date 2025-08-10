package appType

import "strings"

type ImageStorage int

const (
	ImageStorageUnknown ImageStorage = iota // 未知存储类型
	ImageStorageLocal                       // 本地存储
	ImageStorageQiniu                       // 七牛云存储
)

type ImageStorageStrings []string

var imageStorageStrings = ImageStorageStrings{
	"unknown",
	"local",
	"qiniu",
}

func (i ImageStorage) String() string {
	if i < 0 || int(i) >= len(imageStorageStrings) {
		return imageStorageStrings[ImageStorageUnknown]
	}
	return imageStorageStrings[i]
}

// ParseImageStorage 将字符串解析为存储类型
func ParseImageStorage(s string) ImageStorage {
	for i, str := range imageStorageStrings {
		if str == strings.ToLower(s) {
			return ImageStorage(i)
		}
	}
	return ImageStorageUnknown
}
