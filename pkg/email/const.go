package email

const (
	EmailVerifyActionTypeBind = iota
	EmailVerifyActionTypeUnbind
	EmailVerifyActionTypeChangePassword
	EmailVerifyActionTypeResetPassword
	EmailVerifyActionTypeRemoveAccount
)

var EmailVerifyActionTypeMap = map[int]string{
	EmailVerifyActionTypeBind:           "绑定邮箱",
	EmailVerifyActionTypeUnbind:         "解绑邮箱",
	EmailVerifyActionTypeChangePassword: "修改密码",
	EmailVerifyActionTypeResetPassword:  "重置密码",
	EmailVerifyActionTypeRemoveAccount:  "删除账户",
}
