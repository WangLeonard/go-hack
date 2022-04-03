package package1

const ConstVar = "const v1"

var Var = ConstVar

var global interface{}

func init() {
	global = ConstVar

	global = Var
	_ = global
}
