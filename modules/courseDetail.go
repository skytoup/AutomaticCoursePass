package modules

// 课程详情
type CourseDetailJson struct {
	Id       int
	Status   string          // active
	Chapters []CourseChapter // 章节数据
	Title    string
}

// 课程章节
type CourseChapter struct {
	Title       string
	VideoId     int    `json:"content_id"`
	ContentType string `json:"content_type"` // Video
	Content     CourseChapterContent
}

// 课程章节的内容
type CourseChapterContent struct {
	Seconds int // 视频秒数
}
