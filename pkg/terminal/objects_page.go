package terminal

import (
	"context"
	"errors"
	"fmt"

	"os"
	"path"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/pkg/s3lib"
)

type ObjectsPage struct {
	*ListPage[s3lib.Object]

	context Context
}

func NewObjectsPage(context Context) *ObjectsPage {
	listPage := NewListPage[s3lib.Object]()
	listPage.AddColumn("Name", func(item s3lib.Object) string {
		return strings.TrimPrefix(aws.ToString(item.Object.Key), context.ObjectKey())
	})
	listPage.AddColumn("Size", func(item s3lib.Object) string {
		if item.IsDirectory() {
			return ""
		}
		return humanizeSize(item.Object.Size)
	})
	listPage.AddColumn("Last Modified", func(item s3lib.Object) string {
		if item.IsDirectory() {
			return ""
		}
		return humanizeTime(item.Object.LastModified)
	})

	page := &ObjectsPage{
		ListPage: listPage,
		context:  context,
	}

	listPage.SetSelectedFunc(func(selected s3lib.Object) {
		p := strings.TrimPrefix(aws.ToString(selected.Object.Key), context.ObjectKey())
		if strings.HasSuffix(p, "/") {
			context.OpenPage(NewObjectsPage(context.WithObjectKey(aws.ToString(selected.Object.Key))))
		} else {
			context.OpenPage(NewObjectPage(context.WithObjectKey(aws.ToString(selected.Object.Key))))
		}
	})

	page.load()

	return page
}

func (b *ObjectsPage) Context() Context {
	return b.context
}

func (b *ObjectsPage) Title() string {
	title := "Objects"
	if b.context.ObjectKey() != "" {
		title += " - " + b.context.ObjectKey()
	}

	return title
}

func (b *ObjectsPage) Hotkeys() map[tcell.EventKey]Hotkey {
	return map[tcell.EventKey]Hotkey{
		EventKey(tcell.KeyRune, 'n', 0): {
			Title:   "New Object",
			Handler: func(event *tcell.EventKey) *tcell.EventKey { b.newObjectForm(); return nil },
		},
		EventKey(tcell.KeyRune, 'v', 0): {
			Title: "View Object",
			Handler: func(event *tcell.EventKey) *tcell.EventKey {
				if obj, ok := b.ListPage.GetSelectedRow(); ok {
					err := viewObject(b.context.WithObjectKey(aws.ToString(obj.Object.Key)))
					if err != nil {
						b.context.SetError(err)
					}
				}
				return nil
			},
		},
		EventKey(tcell.KeyRune, 'e', 0): {
			Title: "Edit Object",
			Handler: func(event *tcell.EventKey) *tcell.EventKey {
				if obj, ok := b.ListPage.GetSelectedRow(); ok {
					err := editObject(b.context.WithObjectKey(aws.ToString(obj.Object.Key)))
					if err != nil {
						b.context.SetError(err)
					}
				}
				return nil
			},
		},
	}
}

func (b *ObjectsPage) load() {
	b.ListPage.ClearRows()
	paginator := b.context.S3Client().ListObjects(context.Background(), b.context.Bucket(), b.context.ObjectKey())
	var objects []s3lib.Object
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			b.context.SetError(err)
			return
		}

		objects = append(objects, page...)
	}

	b.ListPage.AddAll(objects)
}

func humanizeTime(t *time.Time) string {
	if t == nil {
		return ""
	}

	return t.In(time.Local).Format("2006-01-02 15:04:05")
}

func humanizeSize(sizep *int64) string {
	if sizep == nil {
		return ""
	}
	size := *sizep

	units := []string{"B", "KiB", "MiB", "GiB", "TiB", "PiB", "EiB"}
	var i int
	for ; size >= 1024 && i < len(units)-1; i++ {
		size /= 1024
	}
	return fmt.Sprintf("%d %s", size, units[i])
}

func (b *ObjectsPage) newObjectForm() {
	b.context.Modal(func(close func()) tview.Primitive {
		form := tview.NewForm()
		form.AddInputField("Name", "", 20, nil, func(text string) {})
		form.AddButton("Create", func() {
			if err := b.createObject(form); err != nil {
				b.context.SetError(err)
			}
			close()
		})
		form.AddButton("Cancel", func() {
			close()
		})
		return form
	}, "newObject", 40, 10)
}

func (b *ObjectsPage) createObject(form *tview.Form) error {
	name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("object name cannot be empty")
	}

	name = b.context.ObjectKey() + name

	tmpDir, err := os.MkdirTemp("", "s3tool")
	if err != nil {
		return err
	}

	tmpFilePath := tmpDir + "/" + name
	if _, err := os.Stat(tmpFilePath); err == nil {
		os.Remove(tmpFilePath)
	}

	err = os.MkdirAll(path.Dir(tmpFilePath), 0755)
	if err != nil {
		return err
	}

	err = os.WriteFile(tmpFilePath, nil, 0644)
	if err != nil {
		return err
	}

	err = EditFile(b.context, tmpFilePath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(tmpFilePath); os.IsNotExist(err) {
		return nil
	}

	err = b.context.S3Client().UploadFile(context.Background(), b.context.Bucket(), name, tmpFilePath)
	if err != nil {
		return err
	}

	b.load()

	return nil
}
