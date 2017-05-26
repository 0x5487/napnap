package napnap

import "strings"

type Router struct {
	tree *tree
}

type tree struct {
	rootNode *node
}

type kind uint8

type node struct {
	parent    *node
	children  []*node
	kind      kind
	name      string
	pNames    []string
	params    []string
	sortOrder int
	handler   *methodHandler
}

type methodHandler struct {
	connect HandlerFunc
	delete  HandlerFunc
	get     HandlerFunc
	head    HandlerFunc
	options HandlerFunc
	patch   HandlerFunc
	post    HandlerFunc
	put     HandlerFunc
	trace   HandlerFunc
}

const (
	// CONNECT HTTP method
	CONNECT = "CONNECT"
	// DELETE HTTP method
	DELETE = "DELETE"
	// GET HTTP method
	GET = "GET"
	// HEAD HTTP method
	HEAD = "HEAD"
	// OPTIONS HTTP method
	OPTIONS = "OPTIONS"
	// PATCH HTTP method
	PATCH = "PATCH"
	// POST HTTP method
	POST = "POST"
	// PUT HTTP method
	PUT = "PUT"
	// TRACE HTTP method
	TRACE = "TRACE"
)

var (
	notFoundHandler = func(c *Context) {
		_logger.debug("NotFound")
	}
)

const (
	skind kind = iota
	pkind
)

// NewRouter function will create a new router instance
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

// Invoke function is a middleware entry
func (r *Router) Invoke(c *Context, next HandlerFunc) {
	h := r.Find(c.Request.Method, c.Request.URL.Path, c)

	if h == nil {
		next(c)
	} else {
		h(c)
	}
}

// All is a shortcut for adding all methods
func (r *Router) All(path string, handler HandlerFunc) {
	r.Add(GET, path, handler)
	r.Add(POST, path, handler)
	r.Add(PUT, path, handler)
	r.Add(DELETE, path, handler)
	r.Add(PATCH, path, handler)
	r.Add(OPTIONS, path, handler)
	r.Add(HEAD, path, handler)
}

// Get is a shortcut for router.Add("GET", path, handle)
func (r *Router) Get(path string, handler HandlerFunc) {
	r.Add(GET, path, handler)
}

// Post is a shortcut for router.Add("POST", path, handle)
func (r *Router) Post(path string, handler HandlerFunc) {
	r.Add(POST, path, handler)
}

// Put is a shortcut for router.Add("PUT", path, handle)
func (r *Router) Put(path string, handler HandlerFunc) {
	r.Add(PUT, path, handler)
}

// Delete is a shortcut for router.Add("DELETE", path, handle)
func (r *Router) Delete(path string, handler HandlerFunc) {
	r.Add(DELETE, path, handler)
}

// Patch is a shortcut for router.Add("PATCH", path, handle)
func (r *Router) Patch(path string, handler HandlerFunc) {
	r.Add(PATCH, path, handler)
}

// Options is a shortcut for router.Add("OPTIONS", path, handle)
func (r *Router) Options(path string, handler HandlerFunc) {
	r.Add(OPTIONS, path, handler)
}

// Head is a shortcut for router.Add("HEAD", path, handle)
func (r *Router) Head(path string, handler HandlerFunc) {
	r.Add(HEAD, path, handler)
}

// Add function which adding path and handler to router
func (r *Router) Add(method string, path string, handler HandlerFunc) {
	_logger.debug("===Add")
	if len(path) == 0 {
		panic("router: path couldn't be empty")
	}
	if path[0:1] != "/" {
		panic("router: path was invalid")
	}
	if len(path) > 1 {
		path = path[1:]
	}
	_logger.debug("path:" + path)

	currentNode := r.tree.rootNode
	if path == "/" {
		_logger.debug("lastNode_param:")
		_logger.debug("method:" + method)
		currentNode.addHandler(method, handler)
		return
	}

	pathArray := strings.Split(path, "/")
	count := len(pathArray)
	pathParams := []string{}

	for index, element := range pathArray {
		if len(element) == 0 {
			continue
		}

		var childNode *node
		if element[0:1] == ":" {
			// that is parameter node
			pName := element[1:]
			_logger.debug("parameterName:" + pName)
			childNode = currentNode.findChildByKind(pkind)
			if childNode == nil {
				childNode = newNode(pName, pkind)
				currentNode.addChild(childNode)
			}

			isFound := false
			for _, p := range childNode.pNames {
				if p == pName {
					isFound = true
				}
			}

			if isFound == false {
				childNode.pNames = append(childNode.pNames, pName)
				_logger.debug("added_parameter_name:" + pName)
			}

			pathParams = append(pathParams, pName)

		} else {
			// that is static node
			childNode = currentNode.findChildByName(element)
			if childNode == nil {
				childNode = newNode(element, skind)
				currentNode.addChild(childNode)
			}
		}

		// last node in the path
		if count == index+1 {
			childNode.params = pathParams
			_logger.debug("lastNode_param:")
			_logger.debug("method:" + method)
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

// Find returns http handler for specific path
func (r *Router) Find(method string, path string, c *Context) HandlerFunc {
	_logger.debug("===Find")
	_logger.debug("method:" + method)
	_logger.debug("path:" + path)
	if path[0:1] == "/" && len(path) > 1 {
		path = path[1:]
	}

	currentNode := r.tree.rootNode
	if path == "/" {
		return currentNode.findHandler(method)
	}

	pathArray := strings.Split(path, "/")
	count := len(pathArray)

	pathParams := make(map[int][]Param)
	var paramsNum int

	for index, element := range pathArray {
		childNode := currentNode.findChildByName(element)

		if childNode == nil {
			// static node was not found and is looking for parameter node
			childNode = currentNode.findChildByKind(pkind)

			if childNode != nil {
				var newParams []Param
				for _, pName := range childNode.pNames {
					param := Param{Key: pName, Value: element}
					newParams = append(newParams, param)
				}
				pathParams[paramsNum] = newParams
				paramsNum++
			}
		}

		if childNode == nil {
			//return notFoundHandler
			return nil
		}

		// last node in the path
		if count == index+1 {
			myHandler := childNode.findHandler(method)
			if myHandler == nil {
				//return notFoundHandler
				_logger.debug("handler was not found")
				return nil
			}

			var newParams []Param
			for _, pName := range childNode.pNames {
				param := Param{Key: pName, Value: element}
				newParams = append(newParams, param)
			}
			pathParams[paramsNum] = newParams
			paramsNum++

			paramsNum = 0
			//println("params_count:", len(pathParams))
			_logger.debug("lastNode_params_count:", len(childNode.params))
			for _, validParam := range childNode.params {
				for _, p := range pathParams[paramsNum] {
					//println("p_value:", index, p.Key+"&"+p.Value)
					if validParam == p.Key {
						_logger.debug("matched: " + validParam + "," + p.Value)
						c.params = append(c.params, p)
					}
				}
				paramsNum++
			}

			return myHandler
		}

		currentNode = childNode
	}
	//return notFoundHandler
	return nil
}

func newNode(name string, t kind) *node {
	return &node{
		kind:      t,
		name:      name,
		sortOrder: 0,
		handler:   &methodHandler{},
	}
}

func (n *node) addChild(node *node) {
	node.parent = n
	n.children = append(n.children, node)
}

func (n *node) findChildByName(name string) *node {
	var result *node
	for _, element := range n.children {
		if strings.EqualFold(element.name, name) && element.kind == skind {
			result = element
			break
		}
	}
	return result
}

func (n *node) findChildByKind(t kind) *node {
	for _, c := range n.children {
		if c.kind == t {
			return c
		}
	}
	return nil
}

func (n *node) addHandler(method string, h HandlerFunc) {
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
	default:
		panic("method was invalid")
	}
}

func (n *node) findHandler(method string) HandlerFunc {
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
		panic("method was invalid")
	}
}

/*
type (
	//MiddlewareFunc func(c *Context, next HandlerFunc)
	// HTTPErrorHandler is a centralized HTTP error handler.
	HTTPErrorHandler func(error)
)


// Conforms to the http.Handler interface.
func (nap *NapNap) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	c := NewContext(req, w)


		for _, m := range nap.middlewares {
			m(c)
			if c.goNext == false {
				break
			}
		}

	// Execute chain
	h := nap.Router.Find(req.Method, req.URL.Path, c)

	if err := h(c); err != nil {
		//nap.httpErrorHandler(err)
	}
}
*/
