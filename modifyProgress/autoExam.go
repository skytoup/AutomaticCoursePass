package modify

import (
	"autoCourse/config"
	"autoCourse/modules"
	"autoCourse/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	f_1 = `"%d":[`
	f_2 = `"%d"`
	c_1 = `,`
	c_2 = `]`
)

// 自动考试
func autoExam(token, usrId *string, courseId int) (success bool) {
	examList := *getExamList(token, courseId)
	for _, exam := range examList {
		examedCount := 0
		if exam.QuizSubmissions != nil {
			examedCount = len(*exam.QuizSubmissions)
		}

		if examedCount != 0 && config.NeedAgainExam == false { // 不是第一次考试 且 不需要继续考试
			util.LogInfo(*token, exam.Id, "考试已提过，不需要继续考试")
			continue
		} else if exam.Allowed_attempts == examedCount {
			util.LogInfo(*token, exam.Id, "考试已提交次数已满")
			continue
		}

		examQuiz := getExamQuiz(token, courseId, exam.Id)
		if examQuiz == nil {
			util.LogInfo(*token, exam.Id, "获取试卷失败")
			continue
		}

		util.LogInfo(*usrId, ":", *examQuiz.Title)
		questionCount := len(*examQuiz.Questions) - 1

		// 构造答案json
		bufStr := bytes.NewBufferString("")
		for qi, question := range *examQuiz.Questions {
			bufStr.WriteString(fmt.Sprintf(f_1, question.Id))
			isFirstAnswer := true
			for _, answer := range *question.Answers {
				if answer.Correct {
					if isFirstAnswer {
						isFirstAnswer = false
					} else {
						bufStr.WriteString(c_1)
					}
					bufStr.WriteString(fmt.Sprintf(f_2, answer.Id))
				}
			}
			bufStr.WriteString(c_2)
			if questionCount != qi {
				bufStr.WriteString(c_1)
			}
		}
		ans := bufStr.String()
		data := fmt.Sprintf(submitExamData, exam.Id, examQuiz.Id, ans, *examQuiz.CreatedAt, examQuiz.Attempt, exam.Id, courseId, *usrId, *examQuiz.CreatedAt, *examQuiz.CreatedAt)

		if success == false {
			success = submitExamQuiz(token, &data, courseId, exam.Id, examQuiz.Id)
		}

		if config.IsExamFirst { // 是否只修改一次
			return
		}
	}
	return
}

// 提交考试答案
func submitExamQuiz(token, data *string, courseId, examId, examQuizId int) (ok bool) {
	msg := "提交答案"

	url := fmt.Sprintf(examSubmitUrl, courseId, examId, examQuizId)

	r, err := createRequest(&url, "PUT")
	if util.LogErrAndStr(msg, err) {
		return
	}

	r.Header.Del("Upgrade-Insecure-Requests")
	r.Header.Del("Cache-Control")
	r.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	r.Header.Set("X-CSRF-Token", csrfToken)
	r.Header.Set("Referer", fmt.Sprintf("http://gdit.class.gaoxiaobang.com/classes/%d", courseId))
	r.Header.Set("X-Requested-With", "XMLHttpRequest")
	r.Header.Set("Origin", "http://gdit.class.gaoxiaobang.com")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")

	r.ContentLength = int64(len(*data))
	r.Body = ioutil.NopCloser(strings.NewReader(*data))

	resp, err := client.Do(r)
	if util.LogErrAndStr(msg, err) {
		return
	}

	if resp.StatusCode != http.StatusOK {
		util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status)
		return
	}

	ok = true
	return
}

// 获取考试内容
func getExamQuiz(token *string, courseId, examId int) (eq *modules.ExamQuiz) {
	msg := "获取考试内容"

	url := fmt.Sprintf(examQuizUrl, courseId, examId)

	r, err := createRequest(&url, "POST")
	if util.LogErrAndStr(msg, err) {
		return
	}

	r.Header.Del("Upgrade-Insecure-Requests")
	r.Header.Del("Cache-Control")
	r.Header.Set("X-CSRF-Token", csrfToken)
	r.Header.Set("Referer", fmt.Sprintf("http://gdit.class.gaoxiaobang.com/classes/%d", courseId))
	r.Header.Set("X-Requested-With", "XMLHttpRequest")
	r.Header.Set("Origin", "http://gdit.class.gaoxiaobang.com")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")

	data := fmt.Sprintf(examQuizData, examId)
	r.ContentLength = int64(len(data))
	r.Body = ioutil.NopCloser(strings.NewReader(data))

	resp, err := client.Do(r)
	if util.LogErrAndStr(msg, err) {
		return
	}

	if resp.StatusCode != http.StatusOK {
		// 500好像是已经回答了两次了
		util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status)
		return
	}

	defer resp.Body.Close()
	jd := json.NewDecoder(resp.Body)

	var m modules.ExamQuiz
	err = jd.Decode(&m)
	if util.LogErr(err) {
		return
	}

	eq = &m
	return
}

// 获取考试列表
func getExamList(token *string, courseId int) (es *[]modules.Exam) {
	msg := "获取考试列表"

	url := fmt.Sprintf(examListUrl, courseId)
	r, err := createRequest(&url, "GET")
	if util.LogErrAndStr(msg, err) {
		return
	}

	r.Header.Del("Upgrade-Insecure-Requests")
	r.Header.Del("Cache-Control")
	r.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	r.Header.Set("X-CSRF-Token", csrfToken)
	r.Header.Set("Referer", fmt.Sprintf("http://gdit.class.gaoxiaobang.com/classes/%d", courseId))
	r.Header.Set("X-Requested-With", "XMLHttpRequest")
	r.Header.Set("Origin", "http://gdit.class.gaoxiaobang.com")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")

	resp, err := client.Do(r)
	if util.LogErrAndStr(msg, err) {
		return
	}

	if resp.StatusCode != http.StatusOK {
		util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status)
		return
	}

	defer resp.Body.Close()
	jd := json.NewDecoder(resp.Body)

	var ms []modules.Exam
	err = jd.Decode(&ms)
	if util.LogErr(err) {
		return
	}

	es = &ms
	return
}
