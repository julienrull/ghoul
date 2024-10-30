package ghoul

type CtxErrorNext struct{}

func (m *CtxErrorNext) Error() string {
	return "There isn't next servable handler."
}
func NewCtxErrorNext() *CtxErrorNext {
    return &CtxErrorNext{}
}
