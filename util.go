package hlsdl

import (
	"encoding/json"
	"fmt"
	"time"
)

func printStruct(v interface{}) {
	d, _ := json.Marshal(v)
	fmt.Println(string(d))
}

func getFilename() string {
	return time.Now().Format("20060102150405") + ".ts"
}
