package test

type MyInter interface {
	GetName(id int) string
}

func GetUser(m MyInter, id int) string {
	user := m.GetName(id)
	return user
}
