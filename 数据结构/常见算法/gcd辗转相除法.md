## 最大公约数简介

最大公约数(GCD/Greatest Common Divisor)指几个整数中共有约数中最大的一个

## 欧几里德算法

欧几里德算法又称辗转相除法, 该算法用于计算两个整数的最大公约数. 定理如下:

```txt
gcd(a,b) = gcd(b,a mod b)
```

```java
// 辗转相除法可以直白一点写hhhh
private int gcd(int a, int b){
       while(b != 0){
           int tmp = b;
           b = a%b;
           a = tmp;
       }
       return a;
   }
```

```javascript
var gcdOfStrings = function(str1, str2) {
  if (str1 + str2 !== str2 + str1) return ''
  const gcd = (a, b) => (0 === b ? a : gcd(b, a % b))
  return str1.substring(0, gcd(str1.length, str2.length))
};

```

[作者](https://leetcode-cn.com/problems/greatest-common-divisor-of-strings/solution/1071-zi-fu-chuan-de-zui-da-gong-yin-zi-by-wonderfu/)

## 更相减损术

```java
class Solution {
    public String gcdOfStrings(String str1, String str2) {
        if (!(str1 + str2).equals(str2 + str1)) {
            return "";
        }
        return gcd(str1,str2);

    }
    //更相减损术
    public String gcd (String s1,String s2){
        if (s1.equals(s2)){
            return s1;
        }
        if (s1.length()>s2.length()){
            s1 = s1.substring(s2.length(), s1.length());
            return gcd(s1,s2);
        }else {
            s2 = s2.substring(s1.length(),s2.length());
            return gcd(s1,s2);
        }
    }
}
```



# 题

[Leetcode](https://leetcode-cn.com/problems/greatest-common-divisor-of-strings)

对于字符串 S 和 T，只有在`S = T + ... + T`（T 与自身连接 1 次或多次）时，我们才认定 “T 能除尽 S”。

返回最长字符串 X，要求满足 X 能除尽 str1 且 X 能除尽 str2。

```txt 
示例 1：

输入：str1 = "ABCABC", str2 = "ABC"
输出："ABC"
示例 2：

输入：str1 = "ABABAB", str2 = "ABAB"
输出："AB"
示例 3：

输入：str1 = "LEET", str2 = "CODE"
输出：""
```

**终止条件**

- 在两个字符串长度相等的时候，进入最终阶段。

  - 如果两个字符串完全相等，则返回该字符串。

  - 如果两个字符串不相等，则返回空字符串。

**缩小参数规模**
如果两个字符串长度不相等，则让str1长度大于str2。
判断str2是否为str1的子集，如果是，则使用str1剩余的部分和str2继续进行比较。
直到两者长度相等，进入终止条件。
```go
func gcdOfStrings(str1 string, str2 string) string {
	l1 := len(str1)
	l2 := len(str2)
	if l1 < l2 {
		return gcdOfStrings(str2, str1)
	}

	if l1 == l2 {
		if str1 == str2 {
			return str1
		}
	
		return ""
	}
	
	for i := 0; i < l2; i++ {
		if str1[i] != str2[i] {
			return ""
		}
	}
	
	return gcdOfStrings(str1[l2:], str2)
}
```

作者：_yunfan_
链接：[LeetCode](https://leetcode-cn.com/problems/greatest-common-divisor-of-strings/solution/godi-gui-jie-fa-by-_yunfan_/)



自己的实现

```go
func gcdOfStrings(str1 string, str2 string) string {
//    if strings.Compare(str1 + str2,str2 + str1) != 0 {
//            return "";
//     }
       if str1 + str2 != str2 + str1 {
            return "";
        }
		s1 := len(str1)
        s2 := len(str2)
        return str1[:gcd(s1,s2)]
}

func gcd(a,b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a % b)
}

```

