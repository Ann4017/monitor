package main

type Monitor struct {
	h C_http
	d C_db
	s C_ses
}

func (m *Monitor) Init(_s_region string, _s_access_key string, _s_secret_key string) error {
	if len(m.d.err_row) != 0 {
		m.s.Init(_s_region, _s_access_key, _s_secret_key)
	}

	return nil
}
