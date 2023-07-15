package r

func Go(f func(val any), val any) {
	go func(val any) {
		defer func() {
			if err := recover(); err != nil {
				// log.Println(err)
			}
		}()

		f(val)
	}(val)
}
