# cookie

- Go

  ```go
  import net/http
  // 查询Cookie
  http.Cookie("key")
  // 设置Cookie
  http.SetCookie(w http.ResponseWriter, cookie *http.Cookie)
  ```

- Gin

  ```go
  c.Cookie("key")
  c.SetCookie("key","value",domain,path,maxAge,secure,httpOnly)
  ```

- 应用场景
    保存HTTP请求的状态:
	
	1.  保存用户登录的状态
	
	2.  保存用户购物车的状态
	
	3.  保存用于定制化的状态
	
	    缺点：
	
	4. 数据最大为4K
	
	5. 保存在browse

# Session

session保存在服务器,依赖于Cookie中的`session ID`

session和cookie的比较:

1. 数据量不受限制cookie最大4k
2. 数据保存在服务端,相对安全

缺点:

后端需要维护一个session服务,可以放到redis,本地etc.会提高系统的复杂度 