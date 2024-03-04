package changes2

type ChangeKeyer struct{}

func (ck ChangeKeyer) GetKey(c *Change) string {
	return c.Key
}
