# 带权路径长度WPL:

树的所有叶结点的带权路径长度之和，称为树的带权路径长度表示为*WPL*

# 最优二叉树/哈夫曼树

WPL最小的二叉树

```c
typedef struct TreeNode *HuffmanTree
    
struct TreeNode {
    int Weight;
    HuffmanTree Left, Right;
}

HuffmanTree Huffman(MinHeap H)
{	// 假设H->Size个权值已经存在H->Elements[]->Weight里
    int i; HuffmanTree T;
    
    BuildMinHeap(H); // 将H->Elements[]按权值调整为最小堆
    
    for (i = 1; i < H->Size; i++) {//做H->Size-1次合并
    	T = malloc( sizeof( struct TreeNode) );
        // 从最小堆种删除一个结点, 作为新T的左子结点
        T->Left = DeleteMin(H);
        
        T->Right = DeleteMin(H);
    	// 计算新权值        
        T->Weight = T->Left->Weight + T->Right->Weight;

        Insert(H, T); // 将新T插入最小堆
    }
    
    T = DeleteMin(H);
    return T;
}
```

$O(N logN)$



# 哈夫曼树的特点

- 没有度为1的结点(一个儿子)
- n个叶子结点的哈夫曼树 共有`2n-1`个结点  // n0 = n2 + 1
- 哈夫曼树的任意非叶节点的左右子树 交换后 仍是哈夫曼树
- 对同一组权值{w1, w2,......., Wn} 存在**不同构的两棵哈夫曼树**



# 哈夫曼编码

- 左边0 右边 1