package rsa

type Rsa interface {
	ToHex(str string, isWrite bool) (string, error)
}
