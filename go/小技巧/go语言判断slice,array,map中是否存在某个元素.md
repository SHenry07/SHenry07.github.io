```go

func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
 
	return false, errors.New("not in array")

————————————————
版权声明：本文为CSDN博主「地质队炼丹师」的原创文章，遵循 CC 4.0 BY-SA 版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/qq_28612967/article/details/97392727
```

```go
package main

import "reflect"

func Contains(array interface{}, val interface{}) (index int) {
    index = -1
    switch reflect.TypeOf(array).Kind() {
        case reflect.Slice: {
            s := reflect.ValueOf(array)
            for i := 0; i < s.Len(); i++ {
                if reflect.DeepEqual(val, s.Index(i).Interface()) {
                    index = i
                    return
                }
            }
        }
    }
    return
}
```

[更详细信息](https://wxnacy.com/2018/11/20/go-in-array/#wgo)