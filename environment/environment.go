package environment

type Env struct {
	Store     map[string]interface{}
	Outer     *Env
}
