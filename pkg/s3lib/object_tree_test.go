package s3lib

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectTree(t *testing.T) {
	tree := NewObjectTree[string]()
	tree.AddObject("folder1/file1.txt", "File 1")
	tree.AddObject("folder1/file2.txt", "File 2")
	tree.AddObject("folder2/file3.txt", "File 3")

	rootItems := tree.ListRootItems()
	if len(rootItems) != 2 {
		t.Errorf("Expected 2 root items, got %d", len(rootItems))
	}

	assert.Equal(t, "folder1/", rootItems[0].Name)
	assert.Equal(t, "", rootItems[0].Item)
	assert.Equal(t, "folder2/file3.txt", rootItems[1].Name)
	assert.Equal(t, "File 3", rootItems[1].Item)
}

func TestSplitObjectName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{"Empty", "", []string{}},
		{"SinglePart", "file.txt", []string{"file.txt"}},
		{"TwoParts", "folder/file.txt", []string{"folder/", "file.txt"}},
		{"MultipleParts", "folder/subfolder/file.txt", []string{"folder/", "subfolder/", "file.txt"}},
		{"LeadingSlash", "/folder/file.txt", []string{"/folder/", "file.txt"}},
		{"TrailingSlash", "folder/file.txt/", []string{"folder/", "file.txt/"}},
		{"MultipleSlashes", "folder//file.txt", []string{"folder//", "file.txt"}},
		{"MixedSlashes", "folder/subfolder//file.txt", []string{"folder/", "subfolder//", "file.txt"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitObjectName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
