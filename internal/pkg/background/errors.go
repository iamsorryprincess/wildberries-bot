package background

type ErrorsChannel struct {
	ch chan error
}

func NewErrorsChannel() *ErrorsChannel {
	return &ErrorsChannel{
		ch: make(chan error, 1),
	}
}

func (c *ErrorsChannel) Push(err error) {
	c.ch <- err
}

func (c *ErrorsChannel) Errors() <-chan error {
	return c.ch
}
