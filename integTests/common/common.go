package common

import (
	"fmt"
	"strconv"
	"time"
)

func GenerateProcessId() string {
	prcId := "test" + strconv.Itoa(int(time.Now().UnixNano()))
	fmt.Println("process id: " + prcId)
	return prcId
}
