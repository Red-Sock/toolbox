package closer

var cl Closer

func Add(f Closable) {
	cl.Add(f)
}

func Close() (err error) {
	return cl.Close()
}
