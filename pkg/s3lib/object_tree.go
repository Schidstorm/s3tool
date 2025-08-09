package s3lib

var delimiter rune = '/'

type objectTree[T any] struct {
	root *objectTreeNode[T]
}

type objectTreeNode[T any] struct {
	Name     string
	Item     T
	children []*objectTreeNode[T]
}

func (t *objectTreeNode[T]) AddChild(keyParts []string, item T) {
	var addNode *objectTreeNode[T]
	if len(keyParts) == 1 {
		addNode = &objectTreeNode[T]{Name: keyParts[0], Item: item}
	} else if len(keyParts) > 1 {
		for _, child := range t.children {
			if child.Name == keyParts[0] {
				child.AddChild(keyParts[1:], item)
				return
			}
		}
		addNode = &objectTreeNode[T]{Name: keyParts[0]}
		addNode.AddChild(keyParts[1:], item)
	} else {
		// This case should not happen, as we always have at least one part.
		return
	}

	t.children = append(t.children, addNode)
}

func (t *objectTreeNode[T]) CombineNameForSimpleNode() (string, T) {
	if len(t.children) == 0 || len(t.children) > 1 {
		return t.Name, t.Item
	}

	suffix, item := t.children[0].CombineNameForSimpleNode()
	return t.Name + suffix, item
}

func NewObjectTree[T any]() *objectTree[T] {
	return &objectTree[T]{}
}

func (t *objectTree[T]) AddObject(objectName string, item T) {
	parts := SplitObjectName(objectName)
	if len(parts) == 0 {
		return // No object name provided.
	}
	if t.root == nil {
		t.root = &objectTreeNode[T]{Name: ""}
	}
	t.root.AddChild(parts, item)
}

type RootItemResult[T any] struct {
	Name string
	Item T
}

func (t *objectTree[T]) ListRootItems() []RootItemResult[T] {
	if t.root == nil {
		return nil
	}

	var items []RootItemResult[T]
	for _, child := range t.root.children {
		name, item := child.CombineNameForSimpleNode()
		items = append(items, RootItemResult[T]{
			Name: name,
			Item: item,
		})
	}
	return items
}

func SplitObjectName(objectName string) []string {
	delimiter := delimiter
	parts := []string{}
	currentPartBegin := 0
	runes := []rune(objectName)
	index := 0
	// add all leading delimiters to the first part
	// kepp all trailing delimiters of current part
	for index < len(runes) && runes[index] == delimiter {
		index++
	}

	for {
		if index >= len(runes) {
			if currentPartBegin < len(runes) {
				part := objectName[currentPartBegin:]
				parts = append(parts, part)
			}
			break
		}

		c := runes[index]

		if c == delimiter {
			for index < len(runes) && runes[index] == delimiter {
				index++
			}

			part := objectName[currentPartBegin:index]
			parts = append(parts, part)
			currentPartBegin = index
		}

		index++
	}

	return parts
}
