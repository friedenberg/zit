package changes

type ChangeKeyer struct{}

func (ck ChangeKeyer) GetKey(c Change) string {
	return c.Key
}

func (ck ChangeKeyer) GetKeyPtr(c *Change) string {
	return c.Key
}
