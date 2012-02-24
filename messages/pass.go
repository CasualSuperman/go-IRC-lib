package irc

const tmpl_pass = "PASS %s\n"

type m_pass struct {
	pass string
}

func NewPassMessage(pass string) m_pass {
	m := m_pass{pass}
	return m
}

func (m m_pass) Tmpl() string {
	return tmpl_pass
}

func (m m_pass) Data() []interface{} {
	return []interface{}{m.pass}
}
