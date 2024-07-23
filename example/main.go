package main

import (
  "fmt"

  "github.com/kuuy/hlsdl"
)

func main() {
  hlsDL := hlsdl.New("https://player.cl9987.com:188/20231216/SW2Ee8Uv/1000kb/hls/index.m3u8", nil, "download", "", 64, true)

  filepath, err := hlsDL.Download()
  if err != nil {
    panic(err)
  }

  fmt.Println(filepath)
}
