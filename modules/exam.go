package modules

type Exam struct {
	Id               int               // 考试id
	Allowed_attempts int               // 可以回答的次数
	QuizSubmissions  *[]QuizSubmission `json:"quiz_submissions"` // 以前提交的答案
}

type QuizSubmission struct {
	Id int
}
