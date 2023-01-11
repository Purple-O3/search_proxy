package objs

type RetreiveDocResp struct {
	Count  int               `json:"count"`
	Result RecallPostingList `json:"result""`
}
