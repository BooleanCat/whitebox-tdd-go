package filething

import "os"

type Remover func(string) error

type FileThing struct {
	Path   string
	remove Remover
}

func New(path string) FileThing {
	return FileThing{
		Path:   path,
		remove: os.Remove,
	}
}

func (fileThing FileThing) Remove() error {
	err := fileThing.remove(fileThing.Path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
