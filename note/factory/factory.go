package factory

// 工廠模式: 案需要設置私人string, 私人method，封裝在公有方法之內來使用，隱藏了某一方的具體步驟

type mes struct { // struct想public用大寫開頭，private則用小寫開頭
	C string
	pwd string
}

func NewMes()*mes {
	return  &mes{}
}

func (m *mes)SetPwd(p string) {
	m.pwd = p
}