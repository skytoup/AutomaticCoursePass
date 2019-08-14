package modules

// 课程列表的数据
type CourseListJson struct {
	ResultJson
	CourseInfo `json:"data"`
}

// 课程信息
type CourseInfo struct {
	Total      int      // 总数
	CourseList []Course `json:"pageList"`
}

// 课程
type Course struct {
	CourseId int `json:"id"`
}
