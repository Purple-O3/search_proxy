package objs

type Posting struct {
	Term     string
	Docid    uint64
	TermFreq int
	Offset   []int
}

type Doc struct {
	Body  string  `json:"body"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

type RecallPosting struct {
	Posting
	Doc
}

type RecallPostingList []RecallPosting

type ResultData struct {
	Repl  RecallPostingList `json:"repl"`
	Docid uint64            `json:"docid"`
}

type RespData struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Result  ResultData `json:"result"`
}
