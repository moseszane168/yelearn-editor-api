package user_test

import (
	"crf-mold/controller/user"
	"fmt"
	"testing"
)

func TestPassWordFormat(t *testing.T) {
	fmt.Println(user.ValidatePassWord("666666"))
}
