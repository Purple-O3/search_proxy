package objs

import "time"

type Posting struct {
	FieldName string `json:"-"`
	Term      string `json:"-"`
	Docid     uint64
	//TermFreq  int
	//Offset    []int
}

type Data struct {
	Modified  string
	Saled     string
	Num       int
	CreatedAt time.Time `search_type:"keyword"`
}

type Doc struct {
	Ident string `search_type:"keyword"`
	Data
}

type RecallPosting struct {
	Posting
	Doc
}

type RecallPostingList []RecallPosting

type RetData struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Count   int               `json:"count"`
	Result  RecallPostingList `json:"result"`
}

// 实现排序
func (h RecallPostingList) Len() int {
	return len(h)
}

func (h RecallPostingList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h RecallPostingList) Less(i, j int) bool {
	if h[i].Docid == h[j].Docid {
		if h[i].FieldName == h[j].FieldName {
			return h[i].Term < h[j].Term
		} else {
			return h[i].FieldName < h[j].FieldName
		}
	} else {
		return h[i].Docid < h[j].Docid
	}
}

type GoRet struct {
	Repl  RecallPostingList
	Count int
	Err   error
}
