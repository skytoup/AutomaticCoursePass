package modules

type ExamQuiz struct {
	Id        int           // 试卷id
	CreatedAt *string       `json:"created_at"` // 创建时间
	Attempt   int           // 第几次提交
	Quiz      `json:"quiz"` // 考试
}

type Quiz struct {
	Title     *string         // 考试标题
	Questions *[]QuizQuestion `json:"quiz_questions"` // 问题列表
}

type QuizQuestion struct {
	Id      int       // 问题id
	Answers *[]Answer // 答案列表
}

type Answer struct {
	Id      int  // 答案id
	Correct bool // 是否正确
}
