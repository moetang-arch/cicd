package setupshared

import "math/rand"

const (
	_STR = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func RandomString(l int) string {
	r := make([]byte, l)
	for i := 0; i < l; i++ {
		r[i] = _STR[rand.Intn(len(_STR))]
	}
	return string(r)
}
