package out

import "io"

type prn struct {
	env PrintEnvironment
}

// NewPrinter creates new Printer interface instance
func NewPrinter(pe PrintEnvironment) Printer {
	return &prn{env: pe}
}

func (r *prn) Cprint(format string, a ...interface{}) {
	r.env.PrintFunc(r.env.Writer(), format, a...)
}

func (r *prn) Println() {
	r.Cprint("\n")
}

func (r *prn) Writer() io.WriteCloser {
	return r.env.Writer()
}
