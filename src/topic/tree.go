package topic

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

type Node struct {
	Children map[string]*Node
	Values   []interface{}
}

func newNode() *Node {
	return &Node{
		Children: make(map[string]*Node),
	}
}

func (n *Node) removeValue(value interface{}) {
	for i, v := range n.Values {
		if v == value {
			// remove without preserving order
			n.Values[i] = n.Values[len(n.Values)-1]
			n.Values = n.Values[:len(n.Values)-1]
			break
		}
	}
}

func (n *Node) clearValues() {
	n.Values = []interface{}{}
}

func (n *Node) string(i int) string {
	str := ""

	if i != 0 {
		str = fmt.Sprintf("%d", len(n.Values))
	}

	for key, node := range n.Children {
		str += fmt.Sprintf("\n| %s'%s' => %s", strings.Repeat(" ", i*2), key, node.string(i+1))
	}

	return str
}

// A Tree implements a thread-safe topic tree.
type Tree struct {
	// The separator character. Default: "/"
	Separator string

	// The single level wildcard character. Default: "+"
	WildcardOne string

	// The multi level wildcard character. Default "#"
	WildcardSome string

	Root  *Node
	mutex sync.RWMutex
}

// NewTree returns a new Tree.
func NewTree() *Tree {
	return &Tree{
		Separator:    "/",
		WildcardOne:  "+",
		WildcardSome: "#",

		Root: newNode(),
	}
}

// Add registers the value for the supplied topic. This function will
// automatically grow the tree. If value already exists for the given topic it
// will not be added again.
func (t *Tree) Add(topic string, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.add(value, 0, strings.Split(topic, t.Separator), t.Root)
}

func (t *Tree) add(value interface{}, i int, segments []string, node *Node) {
	// add value to leaf
	if i == len(segments) {
		for _, v := range node.Values {
			if v == value {
				return
			}
		}

		node.Values = append(node.Values, value)
		return
	}

	segment := segments[i]
	child, ok := node.Children[segment]

	// create missing node
	if !ok {
		child = newNode()
		node.Children[segment] = child
	}

	t.add(value, i+1, segments, child)
}

// Set sets the supplied value as the only value for the supplied topic. This
// function will automatically grow the tree.
func (t *Tree) Set(topic string, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.set(value, 0, strings.Split(topic, t.Separator), t.Root)
}

func (t *Tree) set(value interface{}, i int, segments []string, node *Node) {
	// set value on leaf
	if i == len(segments) {
		node.Values = []interface{}{value}
		return
	}

	segment := segments[i]
	child, ok := node.Children[segment]

	// create missing node
	if !ok {
		child = newNode()
		node.Children[segment] = child
	}

	t.set(value, i+1, segments, child)
}

// Get gets the values from the topic that exactly matches the supplied topics.
func (t *Tree) Get(topic string) []interface{} {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	return t.get(0, strings.Split(topic, t.Separator), t.Root)
}

func (t *Tree) get(i int, segments []string, node *Node) []interface{} {
	// set value on leaf
	if i == len(segments) {
		return node.Values
	}

	// get next segment
	segment := segments[i]
	child, ok := node.Children[segment]
	if !ok {
		return nil
	}

	return t.get(i+1, segments, child)
}

// Remove un-registers the value from the supplied topic. This function will
// automatically shrink the tree.
func (t *Tree) Remove(topic string, value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.remove(value, 0, strings.Split(topic, t.Separator), t.Root)
}

// Empty will unregister all values from the supplied topic. This function will
// automatically shrink the tree.
func (t *Tree) Empty(topic string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.remove(nil, 0, strings.Split(topic, t.Separator), t.Root)
}

func (t *Tree) remove(value interface{}, i int, segments []string, node *Node) bool {
	// clear or remove value from leaf node
	if i == len(segments) {
		if value == nil {
			node.clearValues()
		} else {
			node.removeValue(value)
		}

		return len(node.Values) == 0 && len(node.Children) == 0
	}

	segment := segments[i]
	child, ok := node.Children[segment]

	// node not found
	if !ok {
		return false
	}

	if t.remove(value, i+1, segments, child) {
		delete(node.Children, segment)
	}

	return len(node.Values) == 0 && len(node.Children) == 0
}

// Clear will unregister the supplied value from all topics. This function will
// automatically shrink the tree.
func (t *Tree) Clear(value interface{}) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.clear(value, t.Root)
}

func (t *Tree) clear(value interface{}, node *Node) bool {
	node.removeValue(value)

	// remove value from all nodes
	for segment, child := range node.Children {
		if t.clear(value, child) {
			delete(node.Children, segment)
		}
	}

	return len(node.Values) == 0 && len(node.Children) == 0
}

// Match will return a set of values from topics that match the supplied topic.
// The result set will be cleared from duplicate values.
//
// Note: In contrast to Search, Match does not respect wildcards in the query but
// in the stored tree.
func (t *Tree) Match(topic string) []interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	segments := strings.Split(topic, t.Separator)
	values := t.match([]interface{}{}, 0, segments, t.Root)

	return t.clean(values)
}

func (t *Tree) match(result []interface{}, i int, segments []string, node *Node) []interface{} {
	// add all values to the result set that match multiple levels
	if child, ok := node.Children[t.WildcardSome]; ok {
		result = append(result, child.Values...)
	}

	// when finished add all values to the result set
	if i == len(segments) {
		return append(result, node.Values...)
	}

	// advance children that match a single level
	if child, ok := node.Children[t.WildcardOne]; ok {
		result = t.match(result, i+1, segments, child)
	}

	segment := segments[i]

	// match segments and get children
	if segment != t.WildcardOne && segment != t.WildcardSome {
		if child, ok := node.Children[segment]; ok {
			result = t.match(result, i+1, segments, child)
		}
	}

	return result
}

// MatchFirst will run Match and return the first value or nil.
func (t *Tree) MatchFirst(topic string) interface{} {
	values := t.Match(topic)

	if len(values) > 0 {
		return values[0]
	}

	return nil
}

// Search will return a set of values from topics that match the supplied topic.
// The result set will be cleared from duplicate values.
//
// Note: In contrast to Match, Search respects wildcards in the query but not in
// the stored tree.
func (t *Tree) Search(topic string) []interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	segments := strings.Split(topic, t.Separator)
	values := t.search([]interface{}{}, 0, segments, t.Root)

	return t.clean(values)
}

func (t *Tree) search(result []interface{}, i int, segments []string, node *Node) []interface{} {
	// when finished add all values to the result set
	if i == len(segments) {
		return append(result, node.Values...)
	}

	// get segment
	segment := segments[i]

	// add all current and further values
	if segment == t.WildcardSome {
		result = append(result, node.Values...)

		for _, child := range node.Children {
			result = t.search(result, i, segments, child)
		}
	}

	// add all current values and continue
	if segment == t.WildcardOne {
		result = append(result, node.Values...)

		for _, child := range node.Children {
			result = t.search(result, i+1, segments, child)
		}
	}

	// match segments and get children
	if segment != t.WildcardOne && segment != t.WildcardSome {
		if child, ok := node.Children[segment]; ok {
			result = t.search(result, i+1, segments, child)
		}
	}

	return result
}

// SearchFirst will run Search and return the first value or nil.
func (t *Tree) SearchFirst(topic string) interface{} {
	values := t.Search(topic)

	if len(values) > 0 {
		return values[0]
	}

	return nil
}

// clean will remove duplicates
func (t *Tree) clean(values []interface{}) []interface{} {
	result := values[:0]

	for _, v := range values {
		if contains(result, v) {
			continue
		}

		result = append(result, v)
	}

	return result
}

// Count will count all stored values in the tree. It will not filter out
// duplicate values and thus might return a different result to `len(All())`.
func (t *Tree) Count() int {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.count(0, t.Root)
}

func (t *Tree) count(counter int, node *Node) int {
	// add children to results
	for _, child := range node.Children {
		counter += t.count(counter, child)
	}

	// add values to result
	return counter + len(node.Values)
}

// All will return all stored values in the tree.
func (t *Tree) All() []interface{} {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return t.clean(t.all([]interface{}{}, t.Root))
}

func (t *Tree) all(result []interface{}, node *Node) []interface{} {
	// add children to results
	for _, child := range node.Children {
		result = t.all(result, child)
	}

	// add current node to results
	return append(result, node.Values...)
}

// Reset will completely clear the tree.
func (t *Tree) Reset() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.Root = newNode()
}

// String will return a string representation of the tree.
func (t *Tree) String() string {
	t.mutex.RLock()
	defer t.mutex.RUnlock()

	return fmt.Sprintf("topic.Tree:%s", t.Root.string(0))
}

func contains(list []interface{}, value interface{}) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}

	return false
}

func (t *Tree) Dump() []byte {
	bts, _ := json.Marshal(t.Root)
	return bts
}

func (t *Tree) Recover(bts []byte) error {
	return json.Unmarshal(bts, &t.Root)
}
