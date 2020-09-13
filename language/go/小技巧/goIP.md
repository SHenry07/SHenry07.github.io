```go
package main

import (
   "fmt"
   "net"
   "os"
)

func main() {

   addrs, err := net.InterfaceAddrs()

   if err != nil {
      fmt.Println(err)
      os.Exit(1)
   }

   for _, address := range addrs {
      // 检查ip地址判断是否回环地址
      if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
         if ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
            // fmt.Println(ipnet.IP.IsLinkLocalUnicast())
            // fmt.Println(ipnet.IP.IsInterfaceLocalMulticast())
            // fmt.Println(ipnet.IP.IsLinkLocalMulticast())
            fmt.Println(ipnet.IP.IsGlobalUnicast())
            fmt.Printf("%v\n", ipnet.IP.String())
         }
      }
   }
}
```