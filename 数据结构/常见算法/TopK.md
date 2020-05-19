通常我们需要在一大堆数中求前 $k$大的数。比如在搜索引擎中求当天用户点击次数排名前10000的热词，在文本特征选择中求 $tf-idf$值按从大到小排名前$k$ 等问题，都涉及到一个核心问题，即**TOP-K问题**。

Top-k有两种不同的解法,

1. 使用堆(优先队列)   时间复杂度$O(nlogk)$
2. 类似快速排序的分治法  (平均)时间复杂度$O(n)$

[leetcode例题](https://leetcode-cn.com/problems/zui-xiao-de-kge-shu-lcof)

## 堆

堆的性质是每次可以找出最大或最小的元素。我们可以使用一个大小为 k 的最大堆（大顶堆），将数组中的元素依次入堆，当堆的大小超过 k 时，便将多出的元素从堆顶弹出。我们以数组 $[5, 4, 1, 3, 6, 2, 9],k=3 $为例展示元素入堆的过程，如下面动图所示：

![入堆出堆的过程](https://pic.leetcode-cn.com/8415d7d727b1c78745ac050753531108fa1fa3cbeb0c5a352a506280e7f45c32.gif)

这样，**由于每次从堆顶弹出的数都是堆中最大的，最小的 k 个元素一定会留在堆里**。这样，把数组中的元素全部入堆之后，堆中剩下的 k 个元素就是最大的 k 个数了。

注意在动画中，我们并没有画出堆的内部结构，因为这部分内容并不重要。我们只需要知道堆每次会弹出最大的元素即可。在写代码的时候，我们使用的也是库函数中的优先队列数据结构，如 Java 中的 `PriorityQueue`。在面试中，我们不需要实现堆的内部结构，把数据结构使用好，会分析其复杂度即可。

```go
func getLeastNumbers(arr []int, k int) []int {
    if k <= 0 {
		return nil
	}

	if k > len(arr) {
		return arr
	}

	h := &MaxHeap{}
	*h = append(*h, arr[:k]...)

	heap.Init(h) // 标准库里有

	for i := k; i < len(arr); i++ {
		if arr[i] < (*h)[0] {
			heap.Pop(h)
			heap.Push(h, arr[i])
		}
	}

	return *h
}

type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MaxHeap) Push(x interface{}) {
        *h = append(*h, x.(int))
}
func (h *MaxHeap) Pop() interface{} {
        x := (*h)[len(*h)-1]
        *h = (*h)[:len(*h)-1]
        return x
}
```



## 快排变形

Top K 问题的另一个解法就比较难想到，需要在平时有算法的积累。实际上，“查找第 k 大的元素”是一类算法问题，称为**选择问题**。找第 k 大的数，或者找前 k 大的数，有一个经典的 quick select（快速选择）算法。这个名字和 quick sort（快速排序）看起来很像，算法的思想也和快速排序类似，都是分治法的思想。

让我们回顾快速排序的思路。快速排序中有一步很重要的操作是 partition（划分），从数组中随机选取一个枢纽元素 v，然后原地移动数组中的元素，使得比 v 小的元素在 v 的左边，比 v 大的元素在 v 的右边，如下图所示：

![partiition](https://pic.leetcode-cn.com/7b7adf28123268ce845bed00effcc5d88c8bcc5b9db3e16d0a57909a9e9b07b9.jpg)

这个 partition 操作是原地进行的，需要 O(n)*O*(*n*) 的时间，接下来，快速排序会递归地排序左右两侧的数组。而快速选择（quick select）算法的不同之处在于，接下来只需要递归地选择一侧的数组。快速选择算法想当于一个“不完全”的快速排序，因为我们只需要知道最小的 k 个数是哪些，并不需要知道它们的顺序。

我们的目的是寻找最小的 k*k* 个数。假设经过一次 partition 操作，枢纽元素位于下标 m*m*，也就是说，左侧的数组有 m*m* 个元素，是原数组中最小的 m*m* 个数。那么：

- 若 k = m*k*=*m*，我们就找到了最小的 k*k* 个数，就是左侧的数组；
- 若 k<m*k*<*m* ，则最小的 k*k* 个数一定都在左侧数组中，我们只需要对左侧数组递归地 parition 即可；
- 若 k>m*k*>*m*，则左侧数组中的 m*m* 个数都属于最小的 k*k* 个数，我们还需要在右侧数组中寻找最小的 k-m*k*−*m* 个数，对右侧数组递归地 partition 即可。

这种方法需要多加领会思想，如果你对快速排序掌握得很好，那么稍加推导应该不难掌握 quick select 的要领。

# BFPRT算法原理

BFPRT称为**中位数的中位数算法**，它的最坏时间复杂度为$O(n)$

[复制自知乎](https://zhuanlan.zhihu.com/p/31498036)

**一. 快速排序原理**

先来看看快速排序是如何进行的，一趟快速排序的过程如下

1. 先从序列中选取一个数作为基准数
2. 将比这个数大的数全部放到它的右边，把小于或者等于它的数全部放到它的左边

一趟快速排序也叫做**Partion**，即将序列划分为两部分，一部分比基准数小，另一部分比基准数大，然后再进行分治过程，每一次**Partion**不一定都能保证划分得很均匀，所以最坏情况下的时间复杂度不能保证总是为 ![[公式]](https://www.zhihu.com/equation?tex=O%28n%5Clog%28n%29%29) 。对于**Partion**过程，通常有两种方法

**1. 两个指针从首尾向中间扫描（双向扫描）**

这种方法可以用**挖坑填数**来形容，比如

![img](https://pic3.zhimg.com/80/v2-740c20889c11e2848ed31e03c4da0276_720w.jpg)

初始化：i = 0; j = 9; pivot = a[0];

现在a[0]保存到了变量pivot中了，相当于在数组a[0]处挖了个坑，那么可以将其它的数填到这里来。从j开始向前找一个小于或者等于pivot的数，即将a[8]填入a[0]，但a[8]又形成了一个新坑，再从i开始向后找一个大于pivot的数，即a[3]填入a[8]，那么a[3]又形成了一个新坑......

就这样，直到i==j才停止，最终得到结果如下

![img](https://pic2.zhimg.com/80/v2-5867b2c5d8c5efcad15bb0397f068a49_720w.jpg)

上述过程就是**一趟快速排序**。

```cpp
#include <iostream>
#include <string.h>
#include <stdio.h>
#include <algorithm>
#include <time.h>
 
using namespace std;
const int N = 10005;
 
int Partion(int a[], int l, int r)
{
    int i = l;
    int j = r;
    int pivot = a[l];
    while(i < j)
    {
        while(a[j] >= pivot && i < j)
            j--;
        a[i] = a[j];
        while(a[i] <= pivot && i < j)
            i++;
        a[j] = a[i];
    }
    a[i] = pivot;
    return i;
}
 
void QuickSort(int a[], int l, int r)
{
    if(l < r)
    {
        int k = Partion(a, l, r);
        QuickSort(a, l, k - 1);
        QuickSort(a, k + 1, r);
    }
}
 
int a[N];
 
int main()
{
    int n;
    while(cin >> n)
    {
        for(int i = 0; i < n; i++)
            cin >> a[i];
        QuickSort(a, 0, n - 1);
        for(int i = 0; i < n; i++)
            cout << a[i] << " ";
        cout << endl;
    }
    return 0;
}
```

**2. 两个指针一前一后逐步向前扫描（单向扫描）**

```cpp
#include <iostream>
#include <string.h>
#include <stdio.h>
 
using namespace std;
const int N = 10005;
 
int Partion(int a[], int l, int r)
{
    int i = l - 1;
    int pivot = a[r];
    for(int j = l; j < r; j++)
    {
        if(a[j] <= pivot)
        {
            i++;
            swap(a[i], a[j]);
        }
    }
    swap(a[i + 1], a[r]);
    return i + 1;
}
 
void QuickSort(int a[], int l, int r)
{
    if(l < r)
    {
        int k = Partion(a, l, r);
        QuickSort(a, l, k - 1);
        QuickSort(a, k + 1, r);
    }
}
 
int a[N];
 
int main()
{
    int n;
    while(cin >> n)
    {
        for(int i = 0; i < n; i++)
            cin >> a[i];
        QuickSort(a, 0, n - 1);
        for(int i = 0; i < n; i++)
            cout << a[i] << " ";
        cout << endl;
    }
    return 0;
}
```

基于双向扫描的快速排序要比基于单向扫描的快速排序算法快很多。

**二. BFPRT算法原理**

在BFPTR算法中，仅仅是改变了快速排序**Partion**中的**pivot**值的选取，在快速排序中，我们始终选择第一个元素或者最后一个元素作为**pivot**，而在BFPTR算法中，每次选择五分中位数的中位数作为**pivot**，这样做的目的就是使得划分比较合理，从而避免了最坏情况的发生。算法步骤如下

> **1. 将 ![[公式]](https://www.zhihu.com/equation?tex=n) 个元素划为 ![[公式]](https://www.zhihu.com/equation?tex=%5Clfloor+n%2F5%5Crfloor) 组，每组5个，至多只有一组由 ![[公式]](https://www.zhihu.com/equation?tex=n%5Cbmod5) 个元素组成。**
> **2. 寻找这 ![[公式]](https://www.zhihu.com/equation?tex=%5Clceil+n%2F5%5Crceil) 个组中每一个组的中位数，这个过程可以用插入排序。**
> **3. 对步骤2中的 ![[公式]](https://www.zhihu.com/equation?tex=%5Clceil+n%2F5%5Crceil) 个中位数，重复步骤1和步骤2，递归下去，直到剩下一个数字。**
> **4. 最终剩下的数字即为pivot，把大于它的数全放左边，小于等于它的数全放右边。**
> **5. 判断pivot的位置与k的大小，有选择的对左边或右边递归。**

求第 ![[公式]](https://www.zhihu.com/equation?tex=k) 大就是求第 ![[公式]](https://www.zhihu.com/equation?tex=n-k%2B1) 小，这两者等价。

```cpp
#include <iostream>
#include <string.h>
#include <stdio.h>
#include <time.h>
#include <algorithm>
 
using namespace std;
const int N = 10005;
 
int a[N];
 
//插入排序
void InsertSort(int a[], int l, int r)
{
    for(int i = l + 1; i <= r; i++)
    {
        if(a[i - 1] > a[i])
        {
            int t = a[i];
            int j = i;
            while(j > l && a[j - 1] > t)
            {
                a[j] = a[j - 1];
                j--;
            }
            a[j] = t;
        }
    }
}
 
//寻找中位数的中位数
int FindMid(int a[], int l, int r)
{
    if(l == r) return l;
    int i = 0;
    int n = 0;
    for(i = l; i < r - 5; i += 5)
    {
        InsertSort(a, i, i + 4);
        n = i - l;
        swap(a[l + n / 5], a[i + 2]);
    }
 
    //处理剩余元素
    int num = r - i + 1;
    if(num > 0)
    {
        InsertSort(a, i, i + num - 1);
        n = i - l;
        swap(a[l + n / 5], a[i + num / 2]);
    }
    n /= 5;
    if(n == l) return l;
    return FindMid(a, l, l + n);
}
 
//进行划分过程
int Partion(int a[], int l, int r, int p)
{
    swap(a[p], a[l]);
    int i = l;
    int j = r;
    int pivot = a[l];
    while(i < j)
    {
        while(a[j] >= pivot && i < j)
            j--;
        a[i] = a[j];
        while(a[i] <= pivot && i < j)
            i++;
        a[j] = a[i];
    }
    a[i] = pivot;
    return i;
}
 
int BFPRT(int a[], int l, int r, int k)
{
    int p = FindMid(a, l, r);    //寻找中位数的中位数
    int i = Partion(a, l, r, p);
 
    int m = i - l + 1;
    if(m == k) return a[i];
    if(m > k)  return BFPRT(a, l, i - 1, k);
    return BFPRT(a, i + 1, r, k - m);
}
 
int main()
{
    int n, k;
    scanf("%d", &n);
    for(int i = 0; i < n; i++)
        scanf("%d", &a[i]);
    scanf("%d", &k);
    printf("The %d th number is : %d\n", k, BFPRT(a, 0, n - 1, k));
    for(int i = 0; i < n; i++)
        printf("%d ", a[i]);
    puts("");
    return 0;
}
 
/**
10
72 6 57 88 60 42 83 73 48 85
5
*/
```