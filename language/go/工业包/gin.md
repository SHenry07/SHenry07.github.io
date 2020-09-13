1. 获得所有请求头

   `c.Request.Header`  本质 `http.Header  type Header map[string][]string`

   原生: 在handler中的*request

2. 获得请求参数

   - `  c.Request.ParseMultipartForm(100)` --> `from-data`  -- > `Content-Type : multipart/form-data`可以接受键值对, 文件

   - `  c.Request.ParseForm()` --> `x-www-form-urlencoded`  -- > `Content-Type : application/x-www-form-urlencoded` ==只==可以接受键值对

   - `c.ContentType()`