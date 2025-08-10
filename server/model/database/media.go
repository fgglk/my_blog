package database

type Media struct {
  BaseModel
  Filename    string `gorm:"size:255;not null" json:"filename"` // 原始文件名
  StoragePath string `gorm:"size:512;not null;uniqueIndex" json:"storage_path"` // 存储路径
  FileSize    int64  `json:"file_size"` // 文件大小(字节)
  FileType    string `gorm:"size:50" json:"file_type"` // MIME类型
  Width       int    `json:"width,omitempty"` // 图片宽度(仅图片类型)
  Height      int    `json:"height,omitempty"` // 图片高度(仅图片类型)
  UserID      uint   `gorm:"index;not null" json:"user_id"` // 上传用户
}