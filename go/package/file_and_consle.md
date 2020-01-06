

# 文件操作

## file.Read()

### 基本使用

Read方法定义如下：

```go
func (f *File) Read(b []byte) (n int, err error)
```

它接收一个字节切片，返回读取的字节数和可能的具体错误，读到文件末尾时会返回`0`和`io.EOF`。

```go
file, err := os.Open("file.go") // For read access.
if err != nil {
    log.Fatal(err)
}
```
The file's data can then be read into a slice of bytes. Read and Write take
their byte counts from the length of the argument slice.

    data := make([]byte, 100)
    count, err := file.Read(data)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("read %d bytes: %q\n", count, data[:count])

Note: The maximum number of concurrent operations on a File may be limited
by the OS or the system. The number should be high, but exceeding it may
degrade performance or cause other issues.

## bufio读取文件

bufio是在file的基础上封装了一层API，支持更多的功能。

```go
func readFromFilebyBufio() {
	fileObj, err := os.Open("../main.go")
	if err != nil {
		log.Fatal(err)
	}
	defer fileObj.Close()
	// 创建一个用来从文件中读取内容的对象
	reader := bufio.NewReader(fileObj)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			fmt.Println("读完了")
			return
		}
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(line)
	}

}
```
## ioutil读取整个文件

`io/ioutil`包的`ReadFile`方法能够读取完整的文件，只需要将文件名作为参数传入。

```go
func readFromioutil() {
	ret, err := ioutil.ReadFile("./main.go")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(ret))
}
```

# 文件写入操作

`os.OpenFile()`函数能够以指定模式打开文件，从而实现文件写入相关功能。

```go
func OpenFile(name string, flag int, perm FileMode) (*File, error) {
	...
}
```

其中：

`name`：要打开的文件名 `flag`：打开文件的模式。 模式有以下几种：

|     模式      |   含义   |
| :-----------: | :------: |
| `os.O_WRONLY` |   只写   |
| `os.O_CREATE` | 创建文件 |
| `os.O_RDONLY` |   只读   |
|  `os.O_RDWR`  |   读写   |
| `os.O_TRUNC`  |   清空   |
| `os.O_APPEND` |   追加   |

`perm`：文件权限，一个八进制数。r（读）04，w（写）02，x（执行）01。

## Write和WriteString

```go
func main() {
	file, err := os.OpenFile("xx.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	str := "hello 沙河"
	file.Write([]byte(str))       //写入字节切片数据
	file.WriteString("hello 小王子") //直接写入字符串数据
}
```

## bufio.NewWriter

```go
func main() {
	file, err := os.OpenFile("xx.txt", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open file failed, err:", err)
		return
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	for i := 0; i < 10; i++ {
		writer.WriteString("hello沙河\n") //将数据先写入缓存
	}
	writer.Flush() //将缓存中的内容写入文件
}
```

## ioutil.WriteFile

```go
func main() {
	str := "hello 沙河"
	err := ioutil.WriteFile("./xx.txt", []byte(str), 0666)
	if err != nil {
		fmt.Println("write file failed, err:", err)
		return
	}
}
```

# 终端操作

### fmt.Scan

###  **从缓冲读取输入** bufio

```go
   inputReader = bufio.NewReader(os.Stdin)    //创建一个读取器，并将其与标准输入绑定。
   fmt.Printf("Please enter some input: ")
   input, err = inputReader.ReadString('\n') //读取器对象提供一个方法 ReadString(delim byte) ，该方法从输入中读取内容，直到碰到 delim 指定的字符，然后将读取到的内容连同 delim 字符一起放到缓冲区。
//对unix:使用“\n”作为定界符，而window使用"\r\n"为定界符
     if err == nil {
         fmt.Printf("The input was: %s", input)     
		}
 }
```



```go
// A Logger represents an active logging object that generates lines of
// output to an io.Writer. Each logging operation makes a single call to
// the Writer's Write method. A Logger can be used simultaneously from
// multiple goroutines; it guarantees to serialize access to the Writer.
type Logger struct {
	mu     sync.Mutex // ensures atomic writes; protects the following fields
	prefix string     // prefix to write at beginning of each line
	flag   int        // properties
	out    io.Writer  // destination for output
	buf    []byte     // for accumulating text to write
}

// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line.
// The flag argument defines the logging properties.
func New(out io.Writer, prefix string, flag int) *Logger {
	return &Logger{out: out, prefix: prefix, flag: flag}
}

// SetOutput sets the output destination for the logger.
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.out = w
}

var std = New(os.Stderr, "", LstdFlags)

//这里节选了2段go-log.go的代码https://golang.org/src/log/log.go

// Output writes the output for a logging event. The string s contains
   144  // the text to print after the prefix specified by the flags of the
   145  // Logger. A newline is appended if the last character of s is not
   146  // already a newline. Calldepth is used to recover the PC and is
   147  // provided for generality, although at the moment on all pre-defined
   148  // paths it will be 2.
   149  func (l *Logger) Output(calldepth int, s string) error {
   150  	now := time.Now() // get this early.
   151  	var file string
   152  	var line int
   153  	l.mu.Lock()
   154  	defer l.mu.Unlock()
   155  	if l.flag&(Lshortfile|Llongfile) != 0 {
   156  		// Release lock while getting caller info - it's expensive.
   157  		l.mu.Unlock()
   158  		var ok bool
   159  		_, file, line, ok = runtime.Caller(calldepth)
   160  		if !ok {
   161  			file = "???"
   162  			line = 0
   163  		}
   164  		l.mu.Lock()
   165  	}
   166  	l.buf = l.buf[:0]
   167  	l.formatHeader(&l.buf, now, file, line)
   168  	l.buf = append(l.buf, s...)
   169  	if len(s) == 0 || s[len(s)-1] != '\n' {
   170  		l.buf = append(l.buf, '\n')
   171  	}
   172  	_, err := l.out.Write(l.buf)
   173  	return err
   174  }
```

### 从命令行读取参数

#### os.Args

```go
// os_args.go
package main

import (
    "fmt"
    "os"
    "strings"
)

func main() {
    who := "Alice "
    if len(os.Args) > 1 {
        who += strings.Join(os.Args[1:], " ")
    }
    fmt.Println("Good Morning", who)
}
```

#### flag解析命令行选项
flag包有一个扩展功能用来解析命令行选项。但是通常被用来替换基本常量，例如，在某些情况下我们希望在命令行给常量一些不一样的值。（参看19章的项目）

```go
// 在flag包中一个Flag被定义成一个含有如下字段的结构体：

type Flag struct {
    Name     string // name as it appears on command line
    Usage    string // help message
    Value    Value  // value as set
    DefValue string // default value (as text); for usage message
}
下面的程序echo.go模拟了Unix的echo功能：

package main

import (
    "flag" // command line option parser
    "os"
)

var NewLine = flag.Bool("n", false, "print newline") // echo -n flag, of type *bool

const (
    Space   = " "
    Newline = "\n"
)

func main() {
    flag.PrintDefaults()
    flag.Parse() // Scans the arg list and sets up flags
    var s string = ""
    for i := 0; i < flag.NArg(); i++ {
        if i > 0 {
            s += " "
            if *NewLine { // -n is parsed, flag becomes true
                s += Newline
            }
        }
        s += flag.Arg(i)
    }
    os.Stdout.WriteString(s)
}
```



# 运行时

```go
runtime.Calller(深度) //0 函数本身 1 调用函数 2 二级调用
func Caller(skip int) (pc uintptr, file string, line int, ok bool) {

//获取函数名
    funcname := runtime.FuncForPC(pc).Name() // 没什么记住就好了
```

