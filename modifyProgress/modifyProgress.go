package modify

import (
	"autoCourse/modules"
	"autoCourse/util"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// 标记课程
func tagCourse(token *string, videoId, courseId int) {
	msg := "标记课程"
	url := fmt.Sprintf(courseTagUrl, courseId, videoId)
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

	data := tagCourseData
	r.ContentLength = int64(len(data))
	reader := strings.NewReader(data)
	r.Body = ioutil.NopCloser(reader)

	resp, err := client.Do(r)
	if util.LogErrAndStr(msg, err) {
		return
	}
	if resp.StatusCode != http.StatusOK {
		util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status)
		return
	}
}

// 修改进度
func modifyProgress(courseId, progressId, uid, videoId, videoLen int, time int64, token *string) bool {
	msg := "修改课程进度"
	r, err := createRequest(&modifyProgressUrl, "PUT")
	if util.LogErrAndStr(msg, err) {
		return false
	}

	r.Header.Del("Upgrade-Insecure-Requests")
	r.Header.Del("Cache-Control")
	r.Header.Set("Origin", "http://gdit.class.gaoxiaobang.com")
	r.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	r.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Referer", fmt.Sprintf("http://gdit.class.gaoxiaobang.com/classes/%d", courseId))

	data := fmt.Sprintf(modifyProgressData, progressId, uid, courseId, videoId, videoLen, videoLen, videoLen, videoLen, videoLen, time, time, *token)

	r.ContentLength = int64(len(data))
	reader := strings.NewReader(data)
	r.Body = ioutil.NopCloser(reader)

	resp, err := client.Do(r)
	if util.LogErrAndStr(msg, err) {
		return false
	}
	// curl 访问没有问题，这里一直500错误，body没有信息
	if resp.StatusCode != http.StatusOK {
		util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status)
		// return
	} else {
		bs, _ := ioutil.ReadAll(resp.Body)
		util.LogInfo(msg, "访问成功", resp, string(bs))
	}

	return true
}

// 获取课程进度
func getCourseProgress(token *string, videoId, courseId int) *modules.CourseProgressJson {
	msg := "获取课程进度"
	url := fmt.Sprintf(courseProgressUrl, *token, videoId)
	r, err := createRequest(&url, "GET")
	if util.LogErrAndStr(msg, err) {
		return nil
	}

	r.Header.Del("Upgrade-Insecure-Requests")
	r.Header.Del("Cache-Control")
	r.Header.Set("Origin", "http://gdit.class.gaoxiaobang.com")
	r.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	r.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	r.Header.Set("Referer", fmt.Sprintf("http://gdit.class.gaoxiaobang.com/classes/%d", courseId))
	r.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(r)
	if util.LogErrAndStr(msg, err) {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status)
		return nil
	}
	defer resp.Body.Close()

	b, err := gzip.NewReader(resp.Body)
	if util.LogErrAndStr(msg, err) {
		return nil
	}

	jd := json.NewDecoder(b)

	var m modules.CourseProgressJson
	err = jd.Decode(&m)
	if util.LogErr(err) {
		return nil
	}

	return &m
}

// 获取课程详情
func getCourseInfo(courseId int, token *string) (courseDetails *[]modules.CourseDetailJson) {
	msg := "获取课程信息:"
	url := fmt.Sprintf(courseInfoUrl, courseId)
	r, err := createRequest(&url, "GET")
	if util.LogErrAndStr(msg, err) {
		return
	}

	r.Header.Del("Upgrade-Insecure-Requests")
	r.Header.Del("Cache-Control")
	r.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	r.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	r.Header.Set("X-CSRF-Token", csrfToken)
	r.Header.Set("Referer", fmt.Sprintf("http://gdit.class.gaoxiaobang.com/classes/%d", courseId))
	r.Header.Set("X-Requested-With", "XMLHttpRequest")
	r.Header.Set("Host", "gdit.class.gaoxiaobang.com")

	resp, err := client.Do(r)
	if util.LogErrAndStr(msg, err) {
		return
	}
	if resp.StatusCode != http.StatusOK {
		util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status, r)
		return
	}
	defer resp.Body.Close()

	jd := json.NewDecoder(resp.Body)

	var ms []modules.CourseDetailJson
	err = jd.Decode(&ms)
	if util.LogErrAndStr(msg, err) {
		return
	}
	courseDetails = &ms

	return
}

func getLmsAndCsrf(courseId int) (ok bool) {
	msg := "获取 Lms And Csrf:"
	url := fmt.Sprintf("http://gdit.class.gaoxiaobang.com/classes/%d", courseId)
	r, err := http.NewRequest("GET", url, nil)
	if util.LogErrAndStr(msg, err) {
		return
	}

	r.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	r.Header.Set("Referer", "http://gdit.gaoxiaobang.com/gxb")
	r.Header.Set("Accept-Language", "en-US,en;q=0.5")
	r.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(r)
	if util.LogErrAndStr(msg, err) {
		return
	}
	if resp.StatusCode != http.StatusOK {
		util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status)
		return
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if util.LogErrAndStr(msg, err) {
		return
	}

	body := string(bs)
	token, err := findCSRF_Token(&body)
	if util.LogErrAndStr(msg, err) {
		return
	}
	csrfToken = *token
	util.LogInfo(msg, *token)

	ok = true
	return
}

// 获取课程列表
func getCourseList(token *string) (courseInfo *modules.CourseInfo) {
	msg := "获取课程列表:"
	url := fmt.Sprintf(courseListUrl, *token)
	r, err := createRequest(&url, "GET")
	if util.LogErrAndStr(msg, err) {
		return
	}

	r.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	r.Header.Set("Origin", "http://gdit.gaoxiaobang.com")
	r.Header.Set("Referer", "http://gdit.gaoxiaobang.com/gxb")

	resp, err := client.Do(r)
	if util.LogErrAndStr(msg, err) {
		return
	}
	if resp.StatusCode != http.StatusOK {
		util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status)
		return
	}
	defer resp.Body.Close()

	b, err := gzip.NewReader(resp.Body)
	if util.LogErrAndStr(msg, err) {
		return
	}

	jd := json.NewDecoder(b)
	var courseModule modules.CourseListJson
	err = jd.Decode(&courseModule)
	if util.LogErr(err) {
		return
	}

	if courseModule.Message != "success" {
		return
	}
	courseInfo = &courseModule.CourseInfo
	return
}
