package interpreter

type Frame struct {
	fn *CompiledFunction
	ip int
	bp int
}

func NewFrame(fn *CompiledFunction, ip, bp int) *Frame {
	return &Frame{fn: fn, ip: ip, bp: bp}
}
