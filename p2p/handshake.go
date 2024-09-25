package p2p

type HandShakeFunc func(any) error

func NOPHandShakeFunc(any) error { return nil }