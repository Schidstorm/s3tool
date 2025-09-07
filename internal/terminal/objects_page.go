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
	"github.com/schidstorm/s3tool/internal/s3lib"
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
		EventKey(tcell.KeyRune, 'd', 0): {
			Title: "Delete Object",
			Handler: func(event *tcell.EventKey) *tcell.EventKey {
				if selected, ok := b.GetSelectedRow(); ok {
					b.context.Modal(ConfirmModal("Are you sure you want to delete the object '"+aws.ToString(selected.Object.Key)+"'?", func() {
						b.deleteObject(selected)
					}))
				}

				return nil
			},
		},
	}
}

func (b *ObjectsPage) deleteObject(object s3lib.Object) {
	err := b.context.S3Client().DeleteObject(
		context.Background(),
		b.context.Bucket(),
		aws.ToString(object.Object.Key),
	)
	if err != nil {
		b.context.SetError(err)
	}
	err = b.Load()
	if err != nil {
		b.context.SetError(err)
	}
}

func (b *ObjectsPage) Load() error {
	b.ListPage.ClearRows()
	paginator := b.context.S3Client().ListObjects(context.Background(), b.context.Bucket(), b.context.ObjectKey())
	var objects []s3lib.Object
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			return err
		}

		objects = append(objects, page...)
	}

	b.ListPage.AddAll(objects)
	return nil
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
		return NewModal().
			SetTitle("New Object").
			AddInput().SetLabel("Name").
			AddButtons([]string{"Create", "Cancel"}).
			SetDoneFunc(func(buttonLabel string, values map[string]string) {
				if buttonLabel == "Create" {
					b.createObject(values)
					close()
				}
			})
	})
}

func (b *ObjectsPage) createObject(values map[string]string) {
	name := values["Name"]
	name = strings.TrimSpace(name)
	if name == "" {
		b.context.SetError(errors.New("object name cannot be empty"))
		return
	}

	name = b.context.ObjectKey() + name

	tmpDir, err := os.MkdirTemp("", "s3tool")
	if err != nil {
		b.context.SetError(err)
		return
	}

	tmpFilePath := tmpDir + "/" + name
	if _, err := os.Stat(tmpFilePath); err == nil {
		os.Remove(tmpFilePath)
	}

	err = os.MkdirAll(path.Dir(tmpFilePath), 0755)
	if err != nil {
		b.context.SetError(err)
		return
	}

	err = os.WriteFile(tmpFilePath, nil, 0644)
	if err != nil {
		b.context.SetError(err)
		return
	}

	err = EditFile(b.context, tmpFilePath)
	if err != nil {
		b.context.SetError(err)
		return
	}

	if _, err := os.Stat(tmpFilePath); os.IsNotExist(err) {
		// File was deleted, do nothing
		return
	}

	err = b.context.S3Client().UploadFile(context.Background(), b.context.Bucket(), name, tmpFilePath)
	if err != nil {
		b.context.SetError(err)
		return
	}

	err = b.Load()
	if err != nil {
		b.context.SetError(err)
		return
	}
}
