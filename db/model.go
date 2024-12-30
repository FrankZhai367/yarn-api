package db

type User struct {
	Openid    string `gorm:"Primarykey;Column:openid"`
	NickName  string `gorm:"Column:nick_name"`
	AvatarUrl string `gorm:"Column:avatar_url"`
}

func (User) TableName() string {
	return "user"
}

type Image struct {
	Id       int    `gorm:"Primarykey;Column:id"`
	Filename string `gorm:"Column:filename"`
	Openid   string `gorm:"Column:openid;index"`
}

func (Image) TableName() string {
	return "image"
}

// func (record LoginRecord) GetCreateTime() string {
// 	return time.Time(record.CreatedAt).Format("2006-01-02 15:04:05")
// }

type Counter struct {
	Id     int    `gorm:"Primarykey;Column:id"`
	Name   string `gorm:"Column:name"`
	Count  int    `gorm:"Column:count"`
	Openid string `gorm:"Column:openid;index"`
}

func (Counter) TableName() string {
	return "counter"
}

type MyCourse struct {
	Id       int    `gorm:"Primarykey;Column:id"`
	Openid   string `gorm:"Column:openid;index"`
	CourseId string `gorm:"Column:course_id;index"`
}

func (MyCourse) TableName() string {
	return "my_course"
}

type Finished struct {
	Id       int    `gorm:"Primarykey;Column:id"`
	Openid   string `gorm:"Column:openid;index"`
	ObjectId string `gorm:"Column:object_id;index"`
}

func (Finished) TableName() string {
	return "finished"
}

type Reward struct {
	Openid        string `gorm:"Primarykey;Column:openid"`
	CrochetCount  int    `gorm:"Column:crochet_count"`
	KnittingCount int    `gorm:"Column:knitting_count"`
	Lv1Count      int    `gorm:"Column:lv1_count"`
	Lv2Count      int    `gorm:"Column:lv2_count"`
	Lv3Count      int    `gorm:"Column:lv3_count"`
	ShareCount    int    `gorm:"Column:share_count"`
}

func (Reward) TableName() string {
	return "reward"
}
