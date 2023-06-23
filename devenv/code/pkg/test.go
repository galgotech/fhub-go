package pkg

import "fmt"

type Functions struct {
}

func (f *Functions) FuncTest(a string) string {
	return fmt.Sprintf("out %s", a)
}
