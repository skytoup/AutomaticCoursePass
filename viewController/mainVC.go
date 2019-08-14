package viewController

import (
	"autoCourse/config"
	"autoCourse/db"
	"autoCourse/modifyProgress"
	"autoCourse/modules"
	"autoCourse/util"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

var (
	h404 = []byte("404")

	temp, _ = template.ParseFiles("view/main.html")

	regUsr, _ = regexp.Compile("^01\\d{2}(13|14)\\d{4}$")
	regPsw, _ = regexp.Compile("^[\x21-\x7E]{6,16}$")
	regSql, _ = regexp.Compile(`(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`)
)

type MainVC struct {
}

func init() {
	// 开启定时检测人数，进行修改
	modify.TickerCheckModify(config.CheckWaitCountTime)
}

func createMainDate(tip *string) (d *modules.MainData) {
	d = &modules.MainData{}
	if tip != nil {
		d.Tip = *tip
	}
	totalCount, waitCount, successCount := db.GetAccountCount(&modules.Account{})
	d.HandledCount = successCount
	d.PendingCount = waitCount
	d.SubmittedCount = totalCount
	return
}

func (mainVC *MainVC) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		mainVC.handleGet(w, r)
	case "POST":
		mainVC.handlePost(w, r)
	default:
		mainVC.handleDefault(w, r)
	}
}

func (*MainVC) handleGet(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		w.WriteHeader(http.StatusNotFound)
		w.Write(h404)
		return
	}
	m := createMainDate(nil)
	temp.Execute(w, m)
	util.LogAcceptF("GET IP: %s \t---\t URL: %s \t---\t Head: %s \t---\t UserAgent: %s ", r.RemoteAddr, r.URL, r.Header, r.UserAgent())
}

func (*MainVC) handlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	usr := strings.TrimSpace(r.Form.Get("Usr"))
	psw := strings.TrimSpace(r.Form.Get("Psw"))
	sysPsw := strings.TrimSpace(r.Form.Get("SystemPsw"))

	var tip string
	if len(usr) == 0 {
		tip = "学号为空"
	} else if len(psw) == 0 {
		tip = "密码为空"
	} else if config.IsCheckSysPsw && len(sysPsw) == 0 {
		tip = "系统密码为空"
	} else if regUsr.MatchString(usr) != true {
		tip = "学号有误:" + usr
	} else if regPsw.MatchString(psw) != true {
		tip = "密码有误(据说是6-16位的英文、数字、字符组合)"
	} else if regSql.MatchString(psw) {
		tip = "你的密码有点问题，要不先去修改修改一下"
	} else if config.IsCheckSysPsw && sysPsw != config.SysPsw {
		tip = "系统密码不正确，好好回去反省吧。"
	} else {
		ac := modules.Account{}
		ac.Usr = &usr
		ac.Psw = &psw
		ac.Status = modules.AccountStatusCreate
		ok, err := db.AddAccount(&ac)
		if ok != 1 {
			tip = "请勿重复添加"
		} else if err != nil {
			tip = "添加失败，请稍后再试"
		}
	}
	util.LogAcceptF("POST IP: %v \t---\t URL: %v \t---\t Head: %v \t---\t UserAgent: %v \t---\t Body:%v", r.RemoteAddr, r.URL, r.Header, r.UserAgent(), r.Form.Encode())
	m := createMainDate(&tip)
	if len(tip) != 0 {
		util.LogInfo(usr, tip)
		temp.Execute(w, m)
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		modify.CheckWaitingCount()
	}
}

func (*MainVC) handleDefault(w http.ResponseWriter, r *http.Request) {
	util.LogAcceptF("%v IP: %v \t---\t URL: %v \t---\t Head: %v \t---\t UserAgent: %v ", r.Method, r.RemoteAddr, r.URL, r.Header, r.UserAgent())
	w.WriteHeader(http.StatusNotFound)
	w.Write(h404)
}
