package jsonrpc2ws

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
