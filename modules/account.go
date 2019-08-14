package modules

import (
	"time"
)

const (
	AccountStatusCreate    = "create"     // 初次创建
	AccountStatusLoginFail = "login fail" // 登陆失败
	AccountStatusSuccess   = "success"    // 成功修改
	AccountStatusWait      = "waiting"    // 等待修改
)

type Account struct {
	Usr       *string   `xorm:"varchar(15) notnull pk index"`
	Psw       *string   `xorm:"varchar(20) notnull"`
	Status    string    `xorm:"varchar(20) notnull index"`
	CreatedAt time.Time `xorm:"created"`
}
