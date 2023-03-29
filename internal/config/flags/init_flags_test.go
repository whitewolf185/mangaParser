package flags

import (
	"fmt"
	"os"
	"testing"
)

func TestTmp(t *testing.T){
	InitServiceFlags()
	fmt.Println(os.Getenv("email_account"))
}