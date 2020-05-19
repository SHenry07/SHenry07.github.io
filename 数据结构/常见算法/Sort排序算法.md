# 排序算法划分

排序算法划分方法有：稳定性，内外排序，时空复杂度

按照**稳定性**划分，稳定排序，如果`a`原本在`b`前面，而`a=b`，排序之后`a`仍然在`b`的前面；而不稳定可能出现在`b`之后。

按照**内外排序**划分，内排序，所有排序操作都在内存中完成；外排序 ：由于数据太大，因此把数据放在磁盘中，而排序通过磁盘和内存的数据传输才能进行；

按照**时空复杂度**划分，时间复杂度是指运行时间，空间复杂度运行完一个程序所需内存的大小。

# 常见排序方法

## 选择排序



- Bubble Sort
- Selection Sort
- Insertion Sort
- Shell Sort
- Merge Sort
- Quick Sort

- [Bucket Sort](https://www.cs.usfca.edu/~galles/visualization/BucketSort.html)
- [Counting Sort](https://www.cs.usfca.edu/~galles/visualization/CountingSort.html)
- [Radix Sort](https://www.cs.usfca.edu/~galles/visualization/RadixSort.html)
- [Heap Sort](https://www.cs.usfca.edu/~galles/visualization/HeapSort.html)

# 冒牌排序Bubble_sort

稳定, 数组没有什么优势, 但是在**链表**排序时会更好实现些, 因为它是两两相邻交换

```go
func BubbleSort(A []int, N int) {
    for P := N-1; P >= 0; P--{
        flag := flase 
        for i := 0; i < P; i ++ { // 一趟冒泡
            if A[i] > A[i+1] {
                A[i], A[i+1] = A[i+1], A[i]
                flag = true
            }
        }
        if !flag { break }
    }
}
// 最好O(n) 最坏 逆序O(N^2)
```

# 插入排序

```go
func InsertSort(A []int, N int){
    var i, P int
    var Tmp int// 数组的类型
    for P = 0; P < N; P ++{
        Tmp = A[P] /* 取出未排序序列中的第一个元素*/
        for i = N; j > 0 && A[i] > Tmp; i--{
             A[i] = A[i-1] /*依次与已排序序列中元素比较并右移*/
        }
        A[i] = Tmp /* 放进合适的位置 */
    }
}
```

同上

# 时间复杂度下界

对于下标 i < j, 如果A[i] < A[j], 则称(i,j)是一对**逆序对(inversion)**

- 交换2个相邻元素正好消去一个逆序对

- 插入排序: T(N,I) = O(N + I)

  如果序列基本有序, 则插入排序简单且高效

==定理: 任意N个不同元素组成的序列平均具有`N(N - 1)/4`个逆序对==

定理: 任何仅以交换相邻两元素来排序的算法, 其平均时间复杂度为$\Omega (N^2)$ 即最好最好的情况也是N^2



# 希尔排序(Shell Sort)

不稳定

```c
void ShellSort( ElementType A[], int N )
{ /* 希尔排序 - 用Sedgewick增量序列 */
     int Si, D, P, i;
     ElementType Tmp;
     /* 这里只列出一小部分增量 */
     int Sedgewick[] = {929, 505, 209, 109, 41, 19, 5, 1, 0};
      
     for ( Si=0; Sedgewick[Si]>=N; Si++ ) 
         ; /* 初始的增量Sedgewick[Si]不能超过待排序列长度 */
 
     for ( D=Sedgewick[Si]; D>0; D=Sedgewick[++Si] )
         for ( P=D; P<N; P++ ) { /* 插入排序*/
             Tmp = A[P];
             for ( i=P; i>=D && A[i-D]>Tmp; i-=D )
                 A[i] = A[i-D];
             A[i] = Tmp;
         }
}
```

```go
func ShellSort(A []int, N int) {
    var Si, D, P, i int
    
    Sedgewick := []int{929, 505, 209, 109, 41, 19, 5, 1, 0}
    
    for Si = 0; Sedgewick[Si] >= N; Si++ {}
    
    for D = Sedgewick[Si]; D > 0; D = Sedgewick[Si] { // 增量
        for P = D; P < N; P++ {
            Tmp := A[p]
            for i = P; i >= D && A[i-D] > Tmp; i -= D{
                A[i] = A[i-D]
            }
            A[i] = Tmp
        }
    }
}
```



# 堆排序

不稳定

## 选择排序

```c
void Selection_Sort ( ElementType A[], int N)
{	int i
    for ( i = 0; i < N; i++) {
        // 从A[i]到A[N-1]中找最小元, 并将其位置赋给MinPosition
        MinPosition = ScanForMin(A, i, N-1);
        // 将未排序部分的最小元,换到有序部分的最后位置
        Swap( A[i], A[MinPosition]);
    }
}
```

利用最大堆, 直接将最大的放到后面

定理: 堆排序处理N个不同元素的随机排序的平均比较次数是$2NlogN - O(Nlog logN)$

虽然堆排序给出最佳平均时间复杂度, 但实际效果不如用Sedgewick增量序列的希尔排序

```txt
如果我们要从全球70多亿人口中找出最富有的100个人，有什么排序算法可以保证不用完全排序就能在中途得到结果吗？
插入排序不行：如果最大富翁最后才出现，那么不到最后一步完成，我们都不敢保证说前面100个就是答案了。
希尔排序本质上是插入的变形，肯定也是不行。
归并呢？因为前后两半的元素在整个排序中不会串边，所以只要后半部分有大富翁，就一定得等到最后一步大合并，才能确保前面排的100个是答案。
而堆排序是唯一可以只用100步就保证得到前100个大富翁的算法！当然，在排序之前先要O(N)的时间去建立最大堆。
所以当我们的问题是要从大量的N个数据中找最大/最小的k个元素时，用堆排序是比较快的，可以在O(N+klogN)时间内得到解 —— 当然k比较小才行。对于这种问题，还有另一种方法是：先把前k个元素调整成最小堆（时间为O(k)）；此后每读入一个元素，首先跟堆顶元素比较，如果没有堆顶大，就直接扔掉了；否则把堆顶元素替换掉，做一次下滤。这样总体最坏复杂度是O(k+Nlogk)。
```

如经典的Topk问题,但是后面又有快排

```c
void Swap( ElementType *a, ElementType *b )
{
     ElementType t = *a; *a = *b; *b = t;
}
  
void PercDown( ElementType A[], int p, int N )
{ /* 改编代码4.24的PercDown( MaxHeap H, int p )    */
  /* 将N个元素的数组中以A[p]为根的子堆调整为最大堆 */
    int Parent, Child;
    ElementType X;
 
    X = A[p]; /* 取出根结点存放的值 */
    for( Parent=p; (Parent*2+1)<N; Parent=Child ) {
        Child = Parent * 2 + 1;
        if( (Child!=N-1) && (A[Child]<A[Child+1]) )
            Child++;  /* Child指向左右子结点的较大者 */
        if( X >= A[Child] ) break; /* 找到了合适位置 */
        else  /* 下滤X */
            A[Parent] = A[Child];
    }
    A[Parent] = X;
}
 
void HeapSort( ElementType A[], int N ) 
{ /* 堆排序 */
     int i;
       
     for ( i=N/2-1; i>=0; i-- )/* 建立最大堆 */
         PercDown( A, i, N );
      
     for ( i=N-1; i>0; i-- ) {
         /* 删除最大堆顶 */
         Swap( &A[0], &A[i] ); /* 见代码7.1 */
         PercDown( A, 0, i );
     }
}
```




