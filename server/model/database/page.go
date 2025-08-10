package database

type Page struct {
	BaseModelWithStatus        // 嵌入带状态的基础模型
	Title               string `gorm:"size:200;not null" json:"title"`
	Slug                string `gorm:"size:255;uniqueIndex;not null" json:"slug"` // 页面URL路径
	Content             string `gorm:"type:longtext;not null" json:"content"`
	Template            string `gorm:"size:100;default:'default'" json:"template"` // 自定义模板
	ShowInNav           bool   `gorm:"default:false" json:"show_in_nav"`           // 是否在导航栏显示
	Sort                int    `gorm:"default:0" json:"sort"`
}

func (Page) TableName() string {
	return "pages"
}
