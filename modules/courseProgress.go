package modules

type CourseProgressJson struct {
	Id       int
	Uid      int   `json:"user_id"`
	ViewTime int   `json:"his_view_time"`
	Time     int64 `json:"created_at"`
}
