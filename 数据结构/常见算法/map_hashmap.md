# 简而言之一种key value的键值映射关系

java 中hashmap

python 中 dist

go 中 map

c 中 <>

[leetcode169](https://leetcode-cn.com/problems/majority-element/)

```go
func majorityElement(nums []int) int {
    counts := make(map[int]int,3)
    currentMaxTimes, currentMostValue := 0, 0
    for _, num := range nums {
        // if _, isExist := ans[nums[i]]; isExist {
        //     ans[nums[i]] ++
        // }else{
        //     ans[nums[i]] = 1 
        // }
        counts[num] ++ 
        if counts[num] > currentMaxTimes {
            currentMaxTimes = counts[num]
            currentMostValue = num 
        }
 
    }
    return currentMostValue
}
```

