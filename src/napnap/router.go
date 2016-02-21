package napnap

import "strings"

type (
	Router struct {
		tree *tree
	}

	tree struct {
		rootNode *node
	}

	node struct {
		parent    *node
		children  []*node
		kind      int
		name      string
		method    string
		sortOrder int
		handler   *methodHandler
	}
	methodHandler struct {
		connect NapNapHandleFunc
		delete  NapNapHandleFunc
		get     NapNapHandleFunc
		head    NapNapHandleFunc
		options NapNapHandleFunc
		patch   NapNapHandleFunc
		post    NapNapHandleFunc
		put     NapNapHandleFunc
		trace   NapNapHandleFunc
	}
)

var (
	notFoundHandler = func(c *Context) error {
		println("NotFound")
		return nil
	}
)

func NewRouter() *Router {
	return &Router{
		tree: &tree{
			rootNode: &node{
				parent:    nil,
				children:  []*node{},
				kind:      0,
				name:      "/",
				sortOrder: 0,
				handler:   &methodHandler{},
			},
		},
	}
}

func (r *Router) Add(method string, path string, handler NapNapHandleFunc) {
	if path[0:1] == "/" {
		path = path[1:]
	} else {
		panic("path was wrong")
	}

	pathArray := strings.Split(path, "/")
	count := len(pathArray)

	if count < 0 {

	}

	currentNode := r.tree.rootNode

	for index, element := range pathArray {
		childNode := currentNode.findChildByName(element)
		if childNode == nil {
			// create a new node
			childNode = newNode(element)
			currentNode.addChild(childNode)
		}

		if count == index+1 {
			println("add_node_name: " + childNode.name)
			childNode.addHandler(method, handler)
		}

		currentNode = childNode
	}

	//test
	/*
		testNode := r.tree.rootNode.children[0].children[0]
		println("test_node_name: " + testNode.name)
		if method == POST {
			testNode.handler.post()
		}*/

}

func (r *Router) Find(method string, path string) NapNapHandleFunc {
	if path[0:1] == "/" {
		path = path[1:]
	}

	pathArray := strings.Split(path, "/")
	count := len(pathArray)

	currentNode := r.tree.rootNode

	for index, element := range pathArray {
		childNode := currentNode.findChildByName(element)
		if childNode == nil {
			return notFoundHandler
		}

		if count == index+1 {
			myHandler := childNode.findHandler(method)
			if myHandler == nil {
				return notFoundHandler
			}
			return myHandler
		}

		currentNode = childNode
	}
	return notFoundHandler
}

func newNode(name string) *node {
	return &node{
		kind:      0,
		name:      name,
		sortOrder: 0,
		handler:   &methodHandler{},
	}
}

func (n *node) addChild(node *node) {
	n.children = append(n.children, node)
}

func (n *node) findChildByName(name string) *node {
	var result *node
	for _, element := range n.children {
		if element.name == name {
			result = element
			break
		}
	}
	return result
}

func (n *node) addHandler(method string, h NapNapHandleFunc) {
	switch method {
	case GET:
		n.handler.get = h
	case POST:
		n.handler.post = h
	case PUT:
		n.handler.put = h
	case DELETE:
		n.handler.delete = h
	case PATCH:
		n.handler.patch = h
	case OPTIONS:
		n.handler.options = h
	case HEAD:
		n.handler.head = h
	case CONNECT:
		n.handler.connect = h
	case TRACE:
		n.handler.trace = h
	}
}

func (n *node) findHandler(method string) NapNapHandleFunc {
	switch method {
	case GET:
		return n.handler.get
	case POST:
		return n.handler.post
	case PUT:
		return n.handler.put
	case DELETE:
		return n.handler.delete
	case PATCH:
		return n.handler.patch
	case OPTIONS:
		return n.handler.options
	case HEAD:
		return n.handler.head
	case CONNECT:
		return n.handler.connect
	case TRACE:
		return n.handler.trace
	default:
		return nil
	}
}
