package main

func main() {
	dsl.funcs.register("test", "This is a test function", []dslParamMeta{}, []dslParamMeta{}, func(a ...any) (any, error) {
		return "test", nil // this function should not survive a restoreState() because it has been created
		// after the last storeState() call
	})
	dsl.shell() // launch the interactive shell to manage the language
}
