package modify

import (
	"autoCourse/config"
	"autoCourse/util"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"
)

const (
	tagCourseData      = `{"learning_progression":{"id":null,"score":null,"join":null,"count":null,"title":null}}`
	modifyProgressData = `{"id":%d,"user_id":%d,"course_id":%d,"video_id":%d,"video_length":%d,"cur_view_time":0,"his_view_time":%d,"cur_finish_clarity":"\u9ad8\u6e05","cur_finish_point":%d,"action_list":[{"time":0,"action":"start"},{"time":0,"action":"view"},{"time":0,"action":"resume"},{"time": %d,"action": "finish"}],"cur_view_period":[],"his_view_period":[{"to":%d,"from":0}],"created_at":%d,"updated_at":%d,"access_token":"%s"}`

	examQuizData   = `{"quiz_submission":{"quiz_id":"%d","status":"unsubmitted","id":null,"submission_data":null,"submission_id":null,"score":null,"kept_score":null,"quiz_data":null,"started_at":null,"ended_at":null,"finished_at":null,"attempt":null,"is_master":null,"quiz":null}}`
	submitExamData = `{"quiz_submission":{"quiz_id":%d,"status":"submitted","id":%d,"submission_data":{%s},"submission_id":null,"score":null,"kept_score":null,"quiz_data":null,"started_at":"%s","ended_at":null,"finished_at":null,"attempt":%d,"is_master":true,"quiz":%d,"context_id":%d,"context_type":"Course","user_id":%s,"quiz_version":null,"position":0,"created_at":"%s","updated_at":"%s"}}`
)

var (
	csrfToken = ""
	regLT, _  = regexp.Compile("(LT-)[^\"]+") // 匹配 lt

	loginPageUrl = "https://passport.gaoxiaobang.com/login?service=http://gdit.gaoxiaobang.com/users/service&tenant_id=75" // 登陆页
	loginUrl     = "https://passport.gaoxiaobang.com/login"                                                                // 登陆
	mainPageUrl  = "http://gdit.gaoxiaobang.com"                                                                           // 主页

	courseListUrl     = "https://w-api.gaoxiaobang.com/me/enrollments?access_token=%s&curPage=1&enrollment_source=free&is_conclude=&pageSize=1000" // 课程列表
	courseInfoUrl     = "http://gdit.class.gaoxiaobang.com/classes/%d/units"                                                                       // 课程详情
	courseTagUrl      = "http://gdit.class.gaoxiaobang.com/classes/%d/chapters/%d/learning_progressions"
	courseProgressUrl = "https://w-api.gaoxiaobang.com/video_logs?access_token=%s&video_id=%d" // 课程进度
	modifyProgressUrl = "https://w-api.gaoxiaobang.com/video_logs"                             // 修改进度

	examListUrl   = "http://gdit.class.gaoxiaobang.com/classes/%d/quizzes"                     // 考试题列表
	examQuizUrl   = "http://gdit.class.gaoxiaobang.com/classes/%d/quizzes/%d/quiz_submissions" // 获取考试题目
	examSubmitUrl = "http://gdit.class.gaoxiaobang.com/classes/%d/quizzes/%d/quiz_submissions/%d"

	client = http.DefaultClient
)

func init() {
	timeout, err := time.ParseDuration(config.ClientTimeout)
	client.Timeout = timeout
	if util.LogErrAndStr("设置client超时:", err) {
		return
	}
}

// 查找str
func findStr(str, s, e *string) *string {
	idx := strings.Index(*str, *s) + len(*s)
	if idx == -1 {
		return nil
	}
	bd := (*str)[idx:]
	idx = strings.Index(bd, *e)
	if idx == -1 {
		return nil
	}
	findStr := bd[:idx]
	return &findStr
}

func findCSRF_Token(body *string) (*string, error) {
	findPrefixStr := `authenticity_token`
	findTokenStr := `<meta content="`
	findTokenStr2 := `"`

	idx := strings.Index(*body, findPrefixStr) + len(findPrefixStr)
	if idx == -1 {
		return nil, errors.New("token not found")
	}
	bd := (*body)[idx:]

	token := findStr(&bd, &findTokenStr, &findTokenStr2)
	if token == nil {
		return nil, errors.New("token not found")
	}
	return token, nil
}

// 查找token
func findTokenAndUserId(body *string) (*string, *string, error) {
	findTokenStr := `access_token: '`
	findTokenStr2 := `'`
	findUserIdStr := `current_user: '{"id":`
	findUserIdStr2 := `,`

	token := findStr(body, &findTokenStr, &findTokenStr2)
	if token == nil {
		return nil, nil, errors.New("token not found")
	}

	uid := findStr(body, &findUserIdStr, &findUserIdStr2)
	if uid == nil {
		return nil, nil, errors.New("user id not found")
	}

	return token, uid, nil
}

// 创建一个请求
func createRequest(url *string, method string) (*http.Request, error) {
	r, err := http.NewRequest(method, *url, nil)

	r.Header.Set("Accept-Encoding", "gzip, deflate")
	r.Header.Set("Accept-Language", "zh-CN,zh;q=0.8")
	r.Header.Set("Upgrade-Insecure-Requests", "1")
	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	r.Header.Set("Connection", "keep-alive")
	r.Header.Set("Cache-Control", "max-age=0")

	return r, err
}
