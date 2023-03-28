package grun

import (
	"strings"
)

// 动态路由的实现，使用trie树
type node struct {
	path string  // 总请求路径  /hello/hh/c
	part string  // 该节点的配对路径 /hh
	parent *node // 父节点 /hello
	children []*node // 子节点 /c
	isSpecial bool // 是否是特别匹配 ':' or '*' 为true
}

// 根据part找到第一个匹配的节点
func (n *node) matchChild(part string) *node {
	for _, ch := range n.children { // 遍历子节点
		if ch.part == part || ch.isSpecial { // 如果路径相同或者是特殊的
			return ch
		}
	}
	return nil
}

// 找到所有匹配的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, ch := range n.children {
		if ch.part == part || ch.isSpecial {
			nodes = append(nodes, ch)
		}
	}
	return nodes
}

// 插入节点
// path: 总路径 parts:分路径 height:插入节点的数目
func (n *node) insert(path string, parts []string, height int) {
	// 全部插入，或者匹配到*
	if len(parts) == height {
		n.path = path
		return
	}
	// 找到第一个要插入节点，观察是否存在
	nowPart := parts[height]
	newNode := n.matchChild(nowPart)
	if newNode == nil {
		newNode = &node{
			part: nowPart,
			parent: n,
			isSpecial: nowPart[0] == ':' || nowPart[0] == '*',
		}
		n.children = append(n.children, newNode)
	}	
	newNode.insert(path, parts ,height + 1)
}

// 查找路径, 返回该路径的头节点
func (n *node) search(parts []string, height int) *node {
	// 找到最后一个节点，或者*匹配符
	if len(parts) == height || strings.HasPrefix(n.part, "*"){
		if n.path == "" {
			return nil
		}
		return n
	}
	// 找到目前路径的所有节点
	nowPart := parts[height]
	chNodes := n.matchChildren(nowPart)
	// 遍历所有节点继续寻找
	for _, ch := range chNodes {
		res := ch.search(parts, height + 1)
		if res != nil {
			return res;
		}
	}
	return nil
}