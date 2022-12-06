package executables

const (
	mountBinary    = "mount"
	mkfsExt4Binary = "mkfs.ext4"
)

type FileSystem struct {
	Mount, Mkfs Executable
}

func NewFileSystem() *FileSystem {
	return &FileSystem{
		Mount: NewExecutable(mountBinary),
		Mkfs:  NewExecutable(mkfsExt4Binary),
	}
}

func (f *FileSystem) MountVolume(device, dir string) error {
	_, err := f.Mount.Execute(device, dir)
	return err
}

func (f *FileSystem) Partition(device string) error {
	_, err := f.Mkfs.Execute(device)
	return err
}
