package qb_utils

/**
	HELPER TO ORGANIZE DIR HIERARCHY
     - ROOT
		- WORKSPACE
			- TMP
*/
// ---------------------------------------------------------------------------------------------------------------------
//	d i r
// ---------------------------------------------------------------------------------------------------------------------

var Dir *DirHelper

type DirHelper struct {
}

func (instance *DirHelper) NewCentral(workPath, tempPath string, createSubTemp bool) (response *DirCentral) {
	response = new(DirCentral)
	response.workPath = workPath
	response.tempPath = tempPath
	response.createSubTemp = createSubTemp
	response.edited = true

	return
}

// ---------------------------------------------------------------------------------------------------------------------
//	i n i t
// ---------------------------------------------------------------------------------------------------------------------

func init() {
	Dir = new(DirHelper)
}

type DirCentral struct {
	workPath      string
	tempPath      string
	createSubTemp bool // tmp is sub dir of workspace

	dirRoot string
	dirWork string
	dirTemp string
	pathLog string

	edited bool
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DirCentral) SetRoot(dir string) *DirCentral {
	instance.edited = true
	instance.setRoot(dir)
	return instance
}

func (instance *DirCentral) SetTemp(dir string) *DirCentral {
	instance.edited = true

	instance.tempPath = dir
	return instance
}

func (instance *DirCentral) SetSubTemp(enabled bool) *DirCentral {
	instance.edited = true

	instance.createSubTemp = enabled
	return instance
}

func (instance *DirCentral) DirRoot() string {
	return instance.dirRoot
}

func (instance *DirCentral) DirTemp() string {
	return instance.dirTemp
}

func (instance *DirCentral) DirWork() string {
	return instance.dirWork
}

func (instance *DirCentral) PathLog() string {
	return instance.pathLog
}

func (instance *DirCentral) GetPath(path string) (response string) {
	response = Paths.Absolutize(path, instance.dirRoot)
	_ = Paths.Mkdir(response)
	return
}

func (instance *DirCentral) GetWorkPath(subPath string) (response string) {
	response = Paths.Absolutize(subPath, instance.dirWork)
	_ = Paths.Mkdir(response)
	return
}

func (instance *DirCentral) GetTempPath(subPath string) (response string) {
	response = Paths.Absolutize(subPath, instance.dirTemp)
	_ = Paths.Mkdir(response)
	return
}

func (instance *DirCentral) Refresh() {
	instance.refresh()
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func (instance *DirCentral) setRoot(dir string) *DirCentral {

	instance.dirRoot = Paths.Absolute(dir)

	return instance
}

func (instance *DirCentral) refresh() {
	if nil != instance && instance.edited {
		instance.edited = false

		if len(instance.dirRoot) == 0 {
			instance.setRoot(Paths.Dir(Paths.GetWorkspacePath()))
		}

		instance.dirWork = Paths.Absolutize(instance.workPath, instance.dirRoot)
		if instance.createSubTemp {
			instance.dirTemp = Paths.Absolutize(instance.tempPath, instance.dirWork)
		} else {
			instance.dirTemp = Paths.Absolutize(instance.tempPath, instance.dirRoot)
		}

		// creates paths
		_ = Paths.Mkdir(instance.dirWork + OS_PATH_SEPARATOR)
		if len(instance.dirTemp) > 0 {
			_ = Paths.Mkdir(instance.dirTemp + OS_PATH_SEPARATOR)
		}

		instance.pathLog = Paths.Concat(instance.dirWork, "logging.log")
	}
}
