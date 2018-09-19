package setupshared_test

import (
	"testing"

	"github.com/moetang-arch/cicd/setup/setupshared"
)

func TestRandomString(t *testing.T) {
	r1 := setupshared.RandomString(20)
	r2 := setupshared.RandomString(20)
	for i := 0; i < 1000000000; i++ {
		if r1 == r2 {
			t.Fatal("two string are equals.", r1, r2)
		}
	}
}
