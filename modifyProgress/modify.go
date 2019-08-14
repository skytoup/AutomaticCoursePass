package modify

import (
	"autoCourse/config"
	"autoCourse/db"
	"autoCourse/modules"
	"autoCourse/util"
	"net/http"
	"sync"
)

// 同步修改
var syncGo struct {
	isGo   bool
	locker sync.Mutex
}

// 尝试批量修改，可能不会启动
func TryModify() {
	go func() {
		syncGo.locker.Lock()
		if syncGo.isGo == false {
			syncGo.isGo = true
			syncGo.locker.Unlock()

			acs := *db.GetWatingAndFailAccount()
			if acs == nil {
				util.LogErrStr("get wating account fail")
				return
			}

			for _, ac := range acs {
				ok, loginCode := modifyAndExam(ac.Usr, ac.Psw)
				if ok {
					ac.Status = modules.AccountStatusSuccess
					util.LogInfo(*ac.Usr, "success modify")
				} else if loginCode == http.StatusUnauthorized {
					// login fail
					ac.Status = modules.AccountStatusLoginFail
				} else {
					// modify fail
					util.LogErrStr(*ac.Usr + " modify fail")
					ac.Status = modules.AccountStatusWait
				}
				if db.UpdateAccount(&ac) != true {
					util.LogErrStr(*ac.Usr + " update status fail")
				}
			}

			syncGo.locker.Lock()
			syncGo.isGo = false
			syncGo.locker.Unlock()
		} else {
			util.LogInfo("try modify fail")
			syncGo.locker.Unlock()
		}
	}()
}

// 修改视频进度 并 考试
func modifyAndExam(usr, psw *string) (ok bool, loginCode int) {
	token, usrId, ok, code := loginWebSite(usr, psw)
	loginCode = code
	if ok == false {
		return
	}

	courseInfo := getCourseList(token)
	if courseInfo == nil {
		util.LogErrStr(*usr + " : get course list fail")
		return
	}

	for i := 0; i < config.MaxModifyCourseCount && i < courseInfo.Total; i++ {
		course := courseInfo.CourseList[i]

		if getLmsAndCsrf(course.CourseId) == false {
			return
		}

		courseDetails := getCourseInfo(course.CourseId, token)
		if courseDetails == nil {
			util.LogErrStr(*usr + " : get course detail fail")
			return
		}

		for _, courseDetail := range *courseDetails {
			if courseDetail.Status == "active" || len(courseDetail.Status) == 0 { // status may be is null
				chapterCount := len(courseDetail.Chapters)
				skip := config.IsSkipFirstVideo
				for i := 0; i < chapterCount; i++ {
					// for i := 1; i < 2; i++ {
					chapter := courseDetail.Chapters[i]
					if chapter.ContentType == "Video" { // 黑科技
						if skip { // 跳过每个章节的第一个视频
							skip = false
							continue
						}
						tagCourse(token, chapter.VideoId, course.CourseId) // 标记视频为已看状态
						progress := getCourseProgress(token, chapter.VideoId, course.CourseId)
						if progress == nil {
							util.LogInfo("获取课程进度失败")
						} else {
							util.LogInfo("view time", progress.ViewTime, "seconds", chapter.Content.Seconds)
						}
						if progress.ViewTime < chapter.Content.Seconds {
							util.LogInfoF("%v will modify course id: %d %s , video id: %d %s", *usr, course.CourseId, courseDetail.Title, chapter.VideoId, chapter.Title)
							if modifyProgress(course.CourseId, progress.Id, progress.Uid, chapter.VideoId, chapter.Content.Seconds, progress.Time, token) == false {
								return
							}
						}
					}
				}
			}
		}

		// 自动考试
		if autoExam(token, usrId, course.CourseId) {
			util.LogInfo(*usrId, "考试成功")
		}
	}
	ok = true
	return
}
