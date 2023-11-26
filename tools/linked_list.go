package tools

var emptyNode = &LinkedNode{}

type Option func(list *LinkedList) *LinkedList

func EnableLock(list *LinkedList) *LinkedList {
	return list
}

type LinkedNode struct {
	data interface{}
	next *LinkedNode
}

func (n *LinkedNode) SetData(data interface{}) {
	n.data = data
}

func (n *LinkedNode) GetData() interface{} {
	return n.data
}

// LinkedList 并发双向列表，头部获取数据，尾部插入数据，增删分两把锁
type LinkedList struct {
	head *LinkedNode
	tail *LinkedNode
}

func NewLinkedList(opts ...Option) *LinkedList {
	list := &LinkedList{}
	list.init()
	for _, opt := range opts {
		opt(list)
	}
	return list
}

func (l *LinkedList) check() {
	if l.head == nil || l.tail == nil {
		l.init()
	}
}

func (l *LinkedList) init() {
	node := &LinkedNode{}
	node.next = nil
	node.data = "init_node"
	l.head = node
	l.tail = node
}

func (l *LinkedList) GetFirstNode() *LinkedNode {
	// 两种情况下为true 1.都为nil 2.都只想初始化node
	if l.head == l.tail {
		return nil
	}
	return l.head.next
}

// Push 往队尾加入数据
func (l *LinkedList) Push(data interface{}) {
	l.check()
	l.pushTail(data)
}

// MPush 批量往队尾加入数据
func (l *LinkedList) MPush(datas ...interface{}) {
	l.check()
	for _, data := range datas {
		l.pushTail(data)
	}
}

func (l *LinkedList) pushTail(data interface{}) {
	node := &LinkedNode{}
	node.data = data
	node.next = nil
	l.tail.next = node
	l.tail = node
}

// Pop 从队头数据弹出1个数据
func (l *LinkedList) Pop() interface{} {
	l.check()
	removeNode := l.popHead()
	return removeNode.data
}

// MPop 从队头数据弹出num个数据
func (l *LinkedList) MPop(num int64) []interface{} {
	datas := make([]interface{}, 0)
	var i int64
	for ; i < num; i++ {
		node := l.popHead()
		if node == emptyNode {
			return datas
		}
		datas = append(datas, node.data)
	}
	return datas
}

func (l *LinkedList) popHead() *LinkedNode {
	dataNode := l.head.next
	if dataNode == nil {
		return emptyNode
	}
	l.head.next = dataNode.next
	dataNode.next = nil

	return dataNode
}
