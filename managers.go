package managers

type Files []File

type File struct {
	Name string
	Path string
	Data string
}

type FilesHandler interface {
	ReadSource() Files
	ReadData(Files) Files
}

type Engine struct {
	FilesHandler    FilesHandler
	PackageManagers PackageManagers
}

func Default(fh FilesHandler) *Engine {
	composer := ComposerDefault()
	return &Engine{
		FilesHandler: fh,
		PackageManagers: PackageManagers{
			composer,
		},
	}
}

func (e *Engine) Run() {
	files := e.FilesHandler.ReadSource()
	e.Filter(files)
	e.Report()

}

type PackageManager interface {
	Match(File)
	Filter(Files)
	Report()
	GetFiles() Files
	SetFiles(Files)
}

type PackageManagers []PackageManager

func (e *Engine) Filter(files Files) {
	for _, pm := range e.PackageManagers {
		pm.Filter(files)
	}
}

func (e *Engine) Report() {
	for _, pm := range e.PackageManagers {
		pm.SetFiles(e.FilesHandler.ReadData(pm.GetFiles()))
		pm.Report()
	}
}
