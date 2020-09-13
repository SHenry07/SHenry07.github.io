# 当前工作目录

```go
// More info on Getwd()
// https://golang.org/src/os/getwd.go
// https://gist.github.com/arxdsilva/4f73d6b89c9eac93d4ac887521121120
import(
	"os" 
	"fmt"
	"log"
	"runtime"
	"strings"
)
func main() {
	dir, err := os.Getwd()
	  if err != nil {
		  log.Fatal(err)
	  }
	fmt.Println(dir)
}
```

# 当前执行文件的目录

```go
package main// 推荐办法

import (
    "fmt"
    "os"
    "path/filepath"
)

func main() {
    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)
    fmt.Println(exPath)
}
```
方法二
```go
import (
    "fmt"
    "log"
    "os"
    "path/filepath"
)

func main() {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            log.Fatal(err)
    }
    fmt.Println(dir)
}
```



# 获取特定文件的信息

```go
type FileInfo interface {
	Name() string       // base name of the file
	Size() int64        // length in bytes for regular files; system-dependent for others
	Mode() FileMode     // file mode bits
	ModTime() time.Time // modification time
	IsDir() bool        // abbreviation for Mode().IsDir()
	Sys() interface{}   // underlying data source (can return nil)
}
func (file *File) Stat() (FileInfo, error) {}
// 代码示例
func main(){
	fileobj , err := os.Open("./logs/noramal.log")
	if err != nil {
		fmt.Printf("opne file failed, err %v",err)
	}
	fileinfo , err := fileobj.Stat()
	if err != nil {
		fmt.Printf("get file info failed, err %v",err)
	}
	fmt.Printf("文件大小是:%vB",fileinfo.Size())
}
```

