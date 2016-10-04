package crash

type TestOptions struct {
	InputFile string
	Debug bool
	Quiet bool
	LogJson bool
	Variables map[string]string
	VariablesFiles []string
}
