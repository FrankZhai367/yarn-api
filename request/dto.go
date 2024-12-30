package request

type SetRewardReq struct {
	CrochetCount  int `json:"crochetCount"`
	KnittingCount int `json:"knittingCount"`
	Lv1Count      int `json:"lv1Count"`
	Lv2Count      int `json:"lv2Count"`
	Lv3Count      int `json:"lv3Count"`
	ShareCount    int `json:"shareCount"`
}

type UpdateCounterReq struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type AddCounterReq struct {
	Name   string `json:"name"`
	Count  int    `json:"count"`
	Openid string `json:"openid"`
}

type UpdateUserReq struct {
	AvatarUrl string `json:"avatarUrl"`
	NickName  string `json:"nickName"`
}
