package boxes

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/schidstorm/s3tool/pkg/s3lib"
)

type ObjectsPage struct {
	*tview.Flex

	list       *ListPage
	statusText *tview.TextView

	client     *s3.Client
	bucketName string
}

func NewObjectsPage(client *s3.Client, bucketName string) *ObjectsPage {
	listPage := NewListPage()
	statusText := tview.NewTextView().SetTextAlign(tview.AlignCenter)

	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexRow)
	flex.AddItem(listPage, 0, 1, true)
	flex.AddItem(statusText, 1, 1, false)

	page := &ObjectsPage{
		Flex:       flex,
		list:       listPage,
		statusText: statusText,
		client:     client,
		bucketName: bucketName,
	}

	listPage.SetSelectedFunc(func(columns []string) {
		page.viewFile(columns[0])
	})

	listPage.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyRune:
			if event.Rune() == 'n' {
				page.newObjectForm()
				return nil
			}
			if event.Rune() == 'e' {
				if i, row := page.list.GetSelectedRow(); i >= 0 {
					err := page.editFile(row.Columns[0])
					if err != nil {
						page.setError(err)
					}
				}
				return nil
			}
		}
		return event
	})

	page.load()

	return page
}

func (b *ObjectsPage) Root() tview.Primitive {
	return b
}

func (b *ObjectsPage) Title() string {
	return "Objects in " + b.bucketName
}

func (b *ObjectsPage) Hotkeys() map[string]string {
	return map[string]string{
		"n": "Create Object",
		"e": "Edit Object",
	}
}

func (b *ObjectsPage) SetSearch(search string) {
	b.list.SetSearch(search)
}

func (b *ObjectsPage) load() {
	b.list.ClearRows()
	b.list.AddRow(Row{
		Header:  true,
		Columns: []string{"Name", "Size", "Last Modified"},
	})

	objects, err := s3lib.ListObjects(b.client, context.Background(), b.bucketName)
	if err != nil {
		b.setError(err)
		return
	}

	rows := make([]Row, len(objects))
	for i, obj := range objects {
		rows[i] = Row{
			Header: false,
			Columns: []string{
				*obj.Key,
				humanizeSize(obj.Size),
				aws.ToTime(obj.LastModified).Format("2006-01-02 15:04:05"),
			},
		}
	}
	b.list.AddRows(rows)
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

func (b *ObjectsPage) setError(err error) {
	b.statusText.SetText("Error: " + err.Error())
	b.statusText.SetTextColor(tcell.ColorRed)
}

func (b *ObjectsPage) newObjectForm() {
	modalName := "newObject"
	form := tview.NewForm()
	form.AddInputField("Name", "", 20, nil, func(text string) {})
	form.AddButton("Create", func() {
		if err := b.createObject(form); err != nil {
			b.setError(err)
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
	if name == "" {
		return errors.New("object name cannot be empty")
	}

	tmpDir, err := os.MkdirTemp("", "s3tool")
	if err != nil {
		return err
	}

	tmpFilePath := tmpDir + "/" + name
	if _, err := os.Stat(tmpFilePath); err == nil {
		os.Remove(tmpFilePath)
	}

	err = EditFile(tmpFilePath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(tmpFilePath); os.IsNotExist(err) {
		return nil
	}

	err = s3lib.UploadFile(b.client, context.Background(), b.bucketName, name, tmpFilePath)
	if err != nil {
		return err
	}

	b.load()

	return nil
}

func (b *ObjectsPage) downloadFileToTmp(objectName string) (string, error) {
	tmpDir, err := os.MkdirTemp("", "s3tool")
	if err != nil {
		return "", err
	}
	tmpFilePath := tmpDir + "/" + objectName

	err = s3lib.DownloadFile(b.client, context.Background(), b.bucketName, objectName, tmpFilePath)
	if err != nil {
		return "", err
	}

	return tmpFilePath, nil

}

func (b *ObjectsPage) editFile(objectName string) error {
	tmpFilePath, err := b.downloadFileToTmp(objectName)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFilePath)

	oldHash, err := fileHash(tmpFilePath)
	if err != nil {
		return err
	}

	err = EditFile(tmpFilePath)
	if err != nil {
		return err
	}

	if _, err := os.Stat(tmpFilePath); os.IsNotExist(err) {
		return errors.New("file does not exist after editing")
	}

	newHash, err := fileHash(tmpFilePath)
	if err != nil {
		return err
	}

	if oldHash != newHash {
		err = s3lib.UploadFile(b.client, context.Background(), b.bucketName, objectName, tmpFilePath)
		if err != nil {
			return err
		}
	}
	b.load()

	return nil
}

func (b *ObjectsPage) viewFile(objectName string) error {
	tmpFilePath, err := b.downloadFileToTmp(objectName)
	if err != nil {
		return err
	}
	defer os.Remove(tmpFilePath)

	return ShowFile(tmpFilePath)
}

func fileHash(filePath string) (string, error) {
	sha256Hash := sha256.New()
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var buffer [4096]byte
	for {
		n, err := file.Read(buffer[:])
		if n > 0 {
			sha256Hash.Write(buffer[:n])
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}
	return string(sha256Hash.Sum(nil)), nil
}
