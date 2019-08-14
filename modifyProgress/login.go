package modify

import (
	"autoCourse/util"
	"compress/gzip"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)

// 登陆网站
func loginWebSite(usr, psw *string) (token, uid *string, ok bool, loginCode int) {
	csrfToken = ""
	var err error
	client.Jar, err = cookiejar.New(nil)
	if util.LogErr(err) {
		return
	}

	lt := getLt()
	if lt == nil || len(*lt) == 0 {
		util.LogErrStr(*usr + " : get lt fail")
		return
	}
	util.LogInfo(*usr, " lt:", *lt)

	code := login(lt, usr, psw)
	if code == http.StatusUnauthorized { // 认证失败
		loginCode = code
		return
	}

	token, uid = getToken()
	if token == nil || len(*token) == 0 {
		util.LogErrStr(*usr + " : get token fail")
		return
	}
	if uid == nil || len(*uid) == 0 {
		util.LogErrStr(*usr + " : get uid fail")
		return
	}
	util.LogInfo(*usr, " token:", *token)
	ok = true

	return
}

// 获取Token
func getToken() (token *string, uid *string) {
	msg := "获取Token:"
	r, err := createRequest(&mainPageUrl, "GET")
	if util.LogErrAndStr(msg, err) {
		return
	}
	r.Header.Set("If-None-Match", "W/\"0f038dd8884e9024302dc6c331209df5\"")
	r.Header.Set("Accept-Encoding", "gzip, deflate, sdch")

	resp, err := client.Do(r)
	if util.LogErrAndStr(msg, err) {
		return
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNotModified {
		util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status)
		return
	}

	defer resp.Body.Close()

	b, err := gzip.NewReader(resp.Body)
	if util.LogErrAndStr(msg, err) {
		return
	}

	bs, err := ioutil.ReadAll(b)
	if util.LogErrAndStr(msg, err) {
		return
	}

	body := string(bs)

	token, uid, err = findTokenAndUserId(&body)
	util.LogErrAndStr(msg, err)

	return
}

// 登陆
func login(lt, name, psw *string) (code int) {
	msg := "登陆:"
	r, err := createRequest(&loginUrl, "POST")
	if util.LogErrAndStr(msg, err) {
		return
	}
	r.Header.Set("Origin", "https://passport.gaoxiaobang.com")
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Referer", "https://passport.gaoxiaobang.com/login?service=http%3A%2F%2Fgdit.gaoxiaobang.com%2Fusers%2Fservice&tenant_id=75")

	r.PostForm = make(url.Values)
	r.PostForm.Set("tenant_id", "75")
	r.PostForm.Set("username", *name)
	r.PostForm.Set("password", *psw)
	r.PostForm.Set("lt", *lt)
	r.PostForm.Set("service", "http%3A%2F%2Fgdit.gaoxiaobang.com%2Fusers%2Fservice")
	vs := url.Values{
		"tenant_id": {"75"},
		"username":  {*name},
		"password":  {*psw},
		"lt":        {*lt},
		"service":   {"http://gdit.gaoxiaobang.com/users/service"},
	}
	vsc := vs.Encode()
	r.ContentLength = int64(len(vsc))

	reader := strings.NewReader(vsc)
	r.Body = ioutil.NopCloser(reader)

	resp, err := client.Do(r)

	if resp == nil || !(resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusOK) {
		if resp != nil {
			util.LogInfo(msg, "访问失败", resp.StatusCode, resp.Status)
		} else {
			util.LogErrAndStr(msg, err)
		}
		return
	}

	url, err := url.Parse("https://passport.gaoxiaobang.com/login")
	if util.LogErr(err) {
		return
	}
	ck_1 := &http.Cookie{Name: "sgsa_id", Value: "gaoxiaobang.com|0", Path: "/", Domain: "passport.gaoxiaobang.com"}
	ck_2 := &http.Cookie{Name: "sgsa_vt_187041_191979", Value: "0", Path: "/", Domain: "passport.gaoxiaobang.com"}

	client.Jar.SetCookies(url, []*http.Cookie{ck_1, ck_2})

	return

}

// 获取lt，用于登陆
func getLt() (lt *string) {
	msg := "获取LT:"
	r, err := createRequest(&loginPageUrl, "GET")
	if util.LogErrAndStr(msg, err) {
		return
	}
	r.Header.Set("Accept-Encoding", "gzip, deflate, sdch")
	r.Header.Set("Referer", "Referer")

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

	bs, err := ioutil.ReadAll(b)
	if util.LogErrAndStr(msg, err) {
		return
	}

	body := string(bs)

	l := regLT.FindString(body)
	lt = &l
	return
}
