package main

import (
	"fmt"
	"math/rand"
	"time"
)

type BtreeNode struct {
	//定义一个B树的节点 数据结构
	//需要以下几个结构
	Leaf  bool  //是否是叶子
	N int  //分支的数量
	keys  []int  //存储类型
	Children  []*BtreeNode //指向自己的多个分支节点

}

//任何一个B树都一个根节点
type Btree struct {
	Root *BtreeNode //根节点
	branch int  //分支
}

//新建一个节点
func NewBtreeNode(n int, branch int , leaf bool) *BtreeNode  {
	return &BtreeNode{
		Leaf:     leaf,
		N:        n,
		keys:     make([]int,branch*2-1),   //n个branch  对应2n  不算上root 2n-1
		Children: make([]*BtreeNode,branch*2),
	}
}

func (B *BtreeNode)Search(key int)(node *BtreeNode,index int)  {
	//搜索一个节点
	i := 0                                  // 5 key
	//找到合适的位置  最后一个key大于    //1 3 4 keys[i]
	for i < B.N && B.keys[i] < key{
		i +=1
	}
	if i < B.N && B.keys[i]==key{
		node,index = B,i  //找到
	}else if B.Leaf == false{
		//如果找不到 到分支上继续找
		node,index = B.Children[i].Search(key)  //递归继续查找
	}
	return
}

func (B *BtreeNode)InsertNotFull (branch int,key int) {
	if B==nil{
		return
	}
	//节点插入数据
	i := B.N  //记录叶子节点的总量
	if B.Leaf{   //判断是否分支结构
		for i > 0 && key < B.keys[i-1]{
			B.keys[i] = B.keys[i-1]
			//从后往前移动  插入排序
			i --
		}
		B.keys[i] = key  //插入数据
		B.N ++   //总量加1
	}else {
		for i > 0 && key < B.keys[i-1]{

			//从后往前移动  插入排序
			i --
		}
		c := B.Children[i]  //找到下标 如果我是叶子
		if c!=nil && c.N ==2 *branch -1{
			B.Split(branch,i) //切割
			if key > B.keys[i]{
				i++
			}
		}
		B.Children[i].InsertNotFull(branch,key) //递归插入到子叶
	}
}

func (parent *BtreeNode)Split(branch int,index int)  {
	//切割数组分层
	//在插入的时候使用
	full := parent.Children[index]
	newnode := NewBtreeNode(branch-1,branch,full.Leaf) //新建一个节点 备份
	for i :=0; i < branch-1 ;i ++{
		newnode.keys[i] = full.keys[i+branch]  //跳过扫描过的节点      15
		newnode.Children[i] = full.Children[i+branch]  //数据移动 12 13     17
	}
	newnode.Children[branch-1] = full.Children[branch*2 -1] //处理最后
	full.N = branch - 1
	//新增一个key到children
	for i := parent.N;i>index;i--{
		parent.Children[i] = parent.Children[i-1]
		parent.keys[i+1] = parent.keys[i]  //从后往前移动
	}
	parent.keys[index] = full.keys[branch-1]
	parent.Children[index+1] = newnode  //插入数据增加总量
	parent.N++
}

func (B *BtreeNode) String() string {
	return fmt.Sprintf("{n=%d,leaf=%v,children=%v}",B.N,B.keys,B.Children)
}

func (tree *Btree)Search (key int)  (n *BtreeNode ,index int)  {
	return tree.Root.Search(key)
}

func (tree *Btree)String() string  {
	return tree.Root.String()  //返回数的字符串
}

func (tree *Btree)Insert(key int)  {
	root := tree.Root
	if root.N == 2*tree.branch-1{
		s := NewBtreeNode(0,tree.branch,false)
		tree.Root = s
		s.Children[0] = root
		s.Split(tree.branch,0)
		root.InsertNotFull(tree.branch,key)
	}else {
		root.InsertNotFull(tree.branch,key)
	}
}
//新建一个B树
func NewBtree(branch int) *Btree  {
	return &Btree{
		Root:   NewBtreeNode(0,branch,true),
		branch: branch,
	}
}

func main()  {
	tree := NewBtree(100000)
	for i :=100000;i>0;i--{
		tree.Insert(rand.Int()%100000)
	}
	fmt.Println(tree.String())
	for i:=0 ;i<10000;i++{
		startTime := time.Now()
		fmt.Println(tree.Search(i))
		fmt.Println("一共用了",time.Since(startTime))
	}

}