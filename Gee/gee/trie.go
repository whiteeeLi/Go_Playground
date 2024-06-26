package gee

import "strings"

// 前缀树的节点定义
type node struct {
	pattern  string  // 待匹配路由，例如 /p/:lang
	part     string  // 路由中的一部分，例如 :lang
	children []*node // 子节点，例如 [doc, tutorial, intro]
	isWild   bool    // 是否精确匹配，part 含有 : 或 * 时为true
}

// 第一个匹配成功的节点，用于插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配成功的节点，用于查找
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

// 使用height标记目前需要匹配的part再parts中的偏移量
func (n *node) insert(pattern string, parts []string, height int) {
	//当height 与 parts长度一致，代表已经完成匹配
	//这里与search中的判断语句是匹配的
	if len(parts) == height {
		n.pattern = pattern
		return
	}
	//当前需要匹配的值
	part := parts[height]
	//在子结点中看看有没有这个值，没有的话就需要插入
	child := n.matchChild(part)
	if child == nil {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	//递归处理后续的节点
	child.insert(pattern, parts, height+1)
}

// 如果获取了node指针不为空，就能判断找到了匹配的路径
func (n *node) search(parts []string, height int) *node {
	//len(parts) == height代表探索到底了
	//strings.HasPrefix(n.part, "*")代表当前的part是随意匹配了
	//这两种情况都可以准备开始结算，看找没找到了
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		//这里以返回地pattern表示是否完成匹配
		//pattern只有终点节点才有
		if n.pattern == "" {
			return nil
		}
		return n
	}
	part := parts[height]
	//在子结点中有哪些匹配这个值
	children := n.matchChildren(part)
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
