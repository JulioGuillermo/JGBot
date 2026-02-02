package virtualfilesaddon

import (
	"JGBot/files"
	"JGBot/js/exec"

	"github.com/fastschema/qjs"
)

type VirtualFile struct {
	name        string
	virtualRoot string
}

func WithVirtualFile(name, virtualRoot string) exec.Option {
	return func(e *exec.Executor) error {
		e.AddAddonObj(&VirtualFile{
			name:        name,
			virtualRoot: virtualRoot,
		})
		return nil
	}
}

func (v *VirtualFile) GetName() string {
	return v.name
}

func (v *VirtualFile) GetJSObj(ctx *qjs.Context) (*qjs.Value, error) {
	return qjs.ToJsValue(ctx, v)
}

func (v *VirtualFile) ReadDir(path string) ([]*files.FileInfo, error) {
	return files.ReadDir(v.virtualRoot, path)
}

func (v *VirtualFile) CreateDir(path string) error {
	return files.CreateDir(v.virtualRoot, path)
}

func (v *VirtualFile) DeleteDir(path string) error {
	return files.DeleteDir(v.virtualRoot, path)
}

func (v *VirtualFile) ReadFile(path string) ([]byte, error) {
	return files.ReadFile(v.virtualRoot, path)
}

func (v *VirtualFile) WriteFile(path string, content []byte) error {
	return files.WriteFile(v.virtualRoot, path, content)
}

func (v *VirtualFile) ReadStrFile(path string) (string, error) {
	bytes, err := v.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (v *VirtualFile) WriteStrFile(path string, content string) error {
	return v.WriteFile(path, []byte(content))
}

func (v *VirtualFile) DeleteFile(path string) error {
	return files.DeleteFile(v.virtualRoot, path)
}

func (v *VirtualFile) Info(path string) (*files.FileInfo, error) {
	return files.GetFileInfo(v.virtualRoot, path)
}

func (v *VirtualFile) Exists(path string) (bool, error) {
	return files.Exists(v.virtualRoot, path)
}

func (v *VirtualFile) Join(paths ...string) string {
	return files.PathJoin(paths...)
}

func (v *VirtualFile) Split(path string) []string {
	return files.PathSplit(path)
}

func (v *VirtualFile) Parent(path string) string {
	return files.PathParent(path)
}

func (v *VirtualFile) Name(path string) string {
	return files.PathName(path)
}
