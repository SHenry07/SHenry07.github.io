解决链表的问题常见的技巧

1. 快慢指针的前进方向相同，且它们步伐的「差」是恒定的，根据这种确定性去解决链表中的一些问题

   确切地说，叫「同步指针」可能更好一些

2. 使用递归函数，避免复杂的更改指针变量指向操作，使得求解问题变得简单

3. 设置「虚拟头结点」，避免对链表第 1 个结点做单独讨论，这个思想在数组里我们见过，叫「哨兵」

4. 为链表编写测试函数，进行调试（在下面的参考代码中有），主要是：

   - 从数组得到一个链表
   - 根据当前结点打印当前结点以及后面的结点。
     这两个方法可以非常方便地帮助我们调试关于链表的程序。

[leetcode实践](https://leetcode-cn.com/problems/middle-of-the-linked-list/solution/kuai-man-zhi-zhen-zhu-yao-zai-yu-diao-shi-by-liwei/)

[链表简介](https://leetcode-cn.com/explore/learn/card/linked-list/)

------

在这里，我们为你提供了一个模板，用于解决链表中的双指针问题。

```java
// Initialize slow & fast pointers
ListNode slow = head;
ListNode fast = head;
/**
 * Change this condition to fit specific problem.
 * Attention: remember to avoid null-pointer error
 **/
while (slow != null && fast != null && fast.next != null) {
    slow = slow.next;           // move slow pointer one step each time
    fast = fast.next.next;      // move fast pointer two steps each time
    if (slow == fast) {         // change this condition to fit specific problem
        return true;
    }
}
return false;   // change return value to fit specific problem 
```

### 提示

------

它与我们在数组中学到的内容类似。但它可能更棘手而且更容易出错。你应该注意以下几点：

**1. 在调用 next 字段之前，始终检查节点是否为空。**

获取空节点的下一个节点将导致空指针错误。例如，在我们运行 `fast = fast.next.next` 之前，需要检查 `fast` 和 `fast.next` 不为空。

**2. 仔细定义循环的结束条件。**

运行几个示例，以确保你的结束条件不会导致无限循环。在定义结束条件时，你必须考虑我们的第一点提示。

### 复杂度分析

------

空间复杂度分析容易。如果只使用指针，而不使用任何其他额外的空间，那么空间复杂度将是 `O(1)`。但是，时间复杂度的分析比较困难。为了得到答案，我们需要分析`运行循环的次数`。

在前面的查找循环示例中，假设我们每次移动较快的指针 2 步，每次移动较慢的指针 1 步。

1. 如果没有循环，快指针需要 `N/2 次`才能到达链表的末尾，其中 N 是链表的长度。
2. 如果存在循环，则快指针需要 `M 次`才能赶上慢指针，其中 M 是列表中循环的长度。

显然，M <= N 。所以我们将循环运行 `N` 次。对于每次循环，我们只需要常量级的时间。因此，该算法的时间复杂度总共为 `O(N)`。

自己分析其他问题以提高分析能力。别忘了考虑不同的条件。如果很难对所有情况进行分析，请考虑最糟糕的情况。