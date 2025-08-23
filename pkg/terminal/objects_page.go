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
	*ListPage

	client     s3lib.Client
	bucketName string
	prefix     string
}

func NewObjectsPage(client s3lib.Client, bucketName, prefix string) *ObjectsPage {
	listPage := NewListPage()
	page := &ObjectsPage{
		ListPage:   listPage,
		client:     client,
		bucketName: bucketName,
		prefix:     prefix,
	}

	listPage.SetSelectedFunc(func(columns []string) {
		if len(columns) < 1 {
			return
		}

		p := columns[0]
		if strings.HasSuffix(p, "/") {
			activeApp.OpenPage(NewObjectsPage(client, bucketName, prefix+p))
		} else {
			activeApp.OpenPage(NewObjectPage(client, bucketName, prefix+p))
		}
	})

	page.load()

	return page
}

func (b *ObjectsPage) Root() tview.Primitive {
	return b
}

func (b *ObjectsPage) Title() string {
	title := "Objects in " + b.bucketName
	if b.prefix != "" {
		title += " - " + b.prefix
	}

	return title
}

func (b *ObjectsPage) Hotkeys() map[tcell.EventKey]Hotkey {
	return map[tcell.EventKey]Hotkey{
		EventKey(tcell.KeyRune, 'n', 0): Hotkey{
			Title:   "New Object",
			Handler: func(event *tcell.EventKey) *tcell.EventKey { b.newObjectForm(); return nil },
		},
		EventKey(tcell.KeyRune, 'v', 0): Hotkey{
			Title: "View Object",
			Handler: func(event *tcell.EventKey) *tcell.EventKey {
				if i, row := b.ListPage.GetSelectedRow(); i >= 0 {
					err := viewObject(b.client, b.bucketName, b.prefix+row.Columns[0])
					if err != nil {
						activeApp.SetError(err)
					}
				}
				return nil
			},
		},
		EventKey(tcell.KeyRune, 'e', 0): Hotkey{
			Title: "Edit Object",
			Handler: func(event *tcell.EventKey) *tcell.EventKey {
				if i, row := b.ListPage.GetSelectedRow(); i >= 0 {
					err := editObject(b.client, b.bucketName, b.prefix+row.Columns[0])
					if err != nil {
						activeApp.SetError(err)
					}
				}
				return nil
			},
		},
	}
}

func (b *ObjectsPage) load() {
	b.ListPage.ClearRows()
	b.ListPage.AddRow(Row{
		Header:  true,
		Columns: []string{"Name", "Size", "Last Modified"},
	})

	paginator := b.client.ListObjects(context.Background(), b.bucketName, b.prefix)
	var objects []s3lib.Object
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.Background())
		if err != nil {
			activeApp.SetError(err)
			return
		}

		objects = append(objects, page...)
	}

	rows := make([]Row, len(objects))
	for i, obj := range objects {
		if obj.IsDirectory() {
			rows[i] = Row{
				Header: false,
				Columns: []string{
					strings.TrimPrefix(aws.ToString(obj.Object.Key), b.prefix),
					"",
					"",
				},
			}
		} else {
			rows[i] = Row{
				Header: false,
				Columns: []string{
					strings.TrimPrefix(aws.ToString(obj.Object.Key), b.prefix),
					humanizeSize(obj.Object.Size),
					humanizeTime(obj.Object.LastModified),
				},
			}
		}
	}
	b.ListPage.AddRows(rows)
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
	modalName := "newObject"
	form := tview.NewForm()
	form.AddInputField("Name", "", 20, nil, func(text string) {})
	form.AddButton("Create", func() {
		if err := b.createObject(form); err != nil {
			activeApp.SetError(err)
		}
		activeApp.CloseModal(modalName)
	})
	form.AddButton("Cancel", func() {
		activeApp.CloseModal(modalName)
	})

	activeApp.Modal(form, modalName, 40, 10)
}

func (b *ObjectsPage) createObject(form *tview.Form) error {
	name := form.GetFormItemByLabel("Name").(*tview.InputField).GetText()
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("object name cannot be empty")
	}

	name = b.prefix + name

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

	err = EditFile(tmpFilePath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(tmpFilePath); os.IsNotExist(err) {
		return nil
	}

	err = b.client.UploadFile(context.Background(), b.bucketName, name, tmpFilePath)
	if err != nil {
		return err
	}

	b.load()

	return nil
}
