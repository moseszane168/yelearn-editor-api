package base

import (
	"database/sql/driver"
	"encoding/json"
)

type FileStruct struct {
	BucketName string `json:"bucketName"` // minio 桶名称
	FileKey    string `json:"fileKey"`    // minio key
	FileName   string `json:"fileName"`   // 文件名称
}

type FileList []*FileStruct

func (f FileList) Value() (driver.Value, error) {
	d, err := json.Marshal(f)
	return string(d), err
}

// 注意，这里的接收器是指针类型，否则无法把数据从数据库读到结构体

func (f *FileList) Scan(v interface{}) error {
	return json.Unmarshal(v.([]byte), f)
}
