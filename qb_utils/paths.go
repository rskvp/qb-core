package qb_utils

import (
	"bufio"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/rskvp/qb-core/qb_rnd"
	"github.com/rskvp/qb-core/qb_sys"
)

type PathsHelper struct {
}

var Paths *PathsHelper

func init() {
	Paths = new(PathsHelper)
}

const DEF_WORKSPACE = "_workspace"

//----------------------------------------------------------------------------------------------------------------------
//	WorkspaceController
//----------------------------------------------------------------------------------------------------------------------

type WorkspaceController struct {
	repo map[string]*Workspace
	mux  sync.Mutex
}

func NewWorkspaceController() *WorkspaceController {
	instance := new(WorkspaceController)
	instance.repo = make(map[string]*Workspace)

	instance.Get("*").SetPath(DEF_WORKSPACE)

	return instance
}

func (instance *WorkspaceController) Get(key string) *Workspace {
	instance.mux.Lock()
	defer instance.mux.Unlock()

	if _, b := instance.repo[key]; !b {
		w := new(Workspace)
		w.name = key
		w.SetPath(DEF_WORKSPACE)
		instance.repo[key] = w
	}
	return instance.repo[key]
}

//----------------------------------------------------------------------------------------------------------------------
//	Workspace
//----------------------------------------------------------------------------------------------------------------------

type Workspace struct {
	name string
	path string
}

func (instance *Workspace) GetPath() string {
	return instance.path
}

func (instance *Workspace) SetPath(path string) {
	instance.path = Paths.Absolute(path)
}

// Resolve get absolute path under this workspace path
func (instance *Workspace) Resolve(partial string) string {
	if filepath.IsAbs(partial) {
		return partial
	}
	return filepath.Join(instance.GetPath(), partial)
}

//----------------------------------------------------------------------------------------------------------------------
//	f i e l d s
//----------------------------------------------------------------------------------------------------------------------

const DEF_TEMP = "./_temp"
const OS_PATH_SEPARATOR = string(os.PathSeparator)

var _workspace *WorkspaceController
var _temp_root string = DEF_TEMP
var _pathSeparator = OS_PATH_SEPARATOR

//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

func init() {
	_workspace = NewWorkspaceController()
	_workspace.Get("*").SetPath(DEF_WORKSPACE)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *PathsHelper) GetWorkspace(name string) *Workspace {
	return _workspace.Get(name)
}

func (instance *PathsHelper) GetWorkspacePath() string {
	return _workspace.Get("*").GetPath()
}

func (instance *PathsHelper) SetWorkspacePath(value string) {
	_workspace.Get("*").SetPath(value)
}

func (instance *PathsHelper) SetWorkspaceParent(value string) {
	_workspace.Get("*").SetPath(filepath.Join(instance.Absolute(value), DEF_WORKSPACE))
}

func (instance *PathsHelper) WorkspacePath(partial string) string {
	return _workspace.Get("*").Resolve(partial)
}

func (instance *PathsHelper) GetTempRoot() string {
	return instance.Absolute(_temp_root)
}

func (instance *PathsHelper) SetTempRoot(path string) {
	_temp_root = instance.Absolute(path)
}

func (instance *PathsHelper) TempPath(partial string) string {
	return instance.Concat(instance.GetTempRoot(), partial)
}

func (instance *PathsHelper) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (instance *PathsHelper) UserHomePath(path string) string {
	dir, err := os.UserHomeDir()
	if nil != err || len(dir) == 0 {
		dir = instance.WorkspacePath(".")
	}
	return instance.Absolutize(path, dir)
}

// DatePath return a date based path ex: "/2020/01/23/file.txt"
// @param root Root, starting directory
// @param partial File name
// @param level 1=Year, 2=Year/Month, 3=Year/Month/Day, 4=Year/Mont/Day/Hour, 5=Year/Mont/Day/Hour/Minute, 6=Year/Mont/Day/Hour/minute/Second
// @param autoCreate If true, creates directories
func (instance *PathsHelper) DatePath(root, partial string, level int, autoCreate bool) string {
	if len(root) == 0 {
		root = "./"
	}
	var pattern string
	switch level {
	case 0:
		pattern = ""
	case 1:
		pattern = "yyyy"
	case 2:
		pattern = "yyyyMM"
	case 3:
		pattern = "yyyyMMdd"
	case 4:
		pattern = "yyyyMMdd/HH"
	case 5:
		pattern = "yyyyMMdd/HH/mm"
	case 6:
		pattern = "yyyyMMdd/HH/mm/ss"
	default:
		pattern = "yyyyMMdd"
	}
	path := Dates.FormatDate(time.Now(), pattern)
	result := instance.Concat(instance.Absolute(root), path, partial)
	if autoCreate {
		_ = instance.Mkdir(result)
	}
	return result
}

func (instance *PathsHelper) Concat(paths ...string) string {
	result := filepath.Join(paths...)
	if instance.IsUrl(result) && strings.Index(result, "://") == -1 {
		result = strings.Replace(result, ":/", "://", 1)
	}
	return result
}

func (instance *PathsHelper) ConcatDir(paths ...string) string {
	result := filepath.Join(paths...) + string(os.PathSeparator)
	if instance.IsUrl(result) && strings.Index(result, "://") == -1 {
		result = strings.Replace(result, ":/", "://", 1)
	}
	return result
}

// Exists Check if a path exists and returns a boolean value or an error if access is denied
// @param path Path to check
func (instance *PathsHelper) Exists(path string) (bool, error) {
	if !strings.Contains(path, "\n") && instance.IsPath(path) {
		_, err := os.Stat(path)
		if err == nil {
			return true, nil
		}
		if os.IsNotExist(err) {
			return false, nil
		}
		return true, err
	}
	return false, nil // not a valid path
}

// EnsureTrailingSlash add path separator to end if any.
func (instance *PathsHelper) EnsureTrailingSlash(dir string) string {
	if strings.HasSuffix(dir, OS_PATH_SEPARATOR) {
		return dir
	}
	if len(instance.Extension(dir)) == 0 {
		return dir + OS_PATH_SEPARATOR
	}
	return dir
}

func (instance *PathsHelper) Absolute(path string) string {
	abs, err := filepath.Abs(path)
	if nil == err {
		return abs
	}
	return path
}

func (instance *PathsHelper) Absolutize(path, root string) string {
	if instance.IsAbs(path) {
		return path
	}
	if len(root) == 0 {
		return instance.Absolute(path)
	}
	return instance.Concat(root, path)
}

func (instance *PathsHelper) Dir(path string) string {
	return filepath.Dir(path)
}

func (instance *PathsHelper) HasExtension(path, ext string) bool {
	return filepath.Ext(path) == ext
}

func (instance *PathsHelper) Extension(path string) string {
	return filepath.Ext(path)
}

func (instance *PathsHelper) ExtensionName(path string) string {
	return strings.Replace(instance.Extension(path), ".", "", 1)
}

func (instance *PathsHelper) FileName(path string, includeExt bool) string {
	if instance.IsUrl(path) {
		uri, err := url.Parse(path)
		if nil != err {
			return ""
		}
		path := uri.Path
		if len(path) > 1 {
			ext := instance.ExtensionName(path)
			if len(ext) > 0 {
				return instance.FileName(path, includeExt)
			}
			return filepath.Base(uri.Path) // i.e. http://file-no-ext
		}
		return ""
	} else {
		base := filepath.Base(path)
		if !includeExt {
			tokens := strings.Split(base, ".")
			if len(tokens) > 2 {
				return strings.Join(tokens[:len(tokens)-1], ".")
			} else {
				ext := filepath.Ext(base)
				return strings.Replace(base, ext, "", 1)
			}
		}
		return base
	}
}

// Mkdir Creates a directory and all subdirectories if does not exists
func (instance *PathsHelper) Mkdir(path string) (err error) {
	// ensure we have a directory
	var abs string
	if filepath.IsAbs(path) {
		abs = path
	} else {
		abs, err = filepath.Abs(path)
	}

	if nil == err {
		if instance.IsFilePath(abs) {
			abs = filepath.Dir(abs)
		}

		if !strings.HasSuffix(abs, _pathSeparator) {
			path = abs + string(os.PathSeparator)
		} else {
			path = abs
		}

		var b bool
		if b, err = instance.Exists(path); !b && nil == err {
			err = os.MkdirAll(path, os.ModePerm)
		}

	}

	return err
}

// IsPath get a string and check if is a valid path
func (instance *PathsHelper) IsPath(path string) bool {
	clean := strings.Trim(path, " ")
	if strings.Contains(clean, OS_PATH_SEPARATOR) {
		return true
	}
	return false
}

// IsTemp return true if path is under "./_temp"
func (instance *PathsHelper) IsTemp(path string) bool {
	tokens := strings.Split(path, _pathSeparator)
	temp := instance.FileName(_temp_root, false)
	for _, token := range tokens {
		if token == temp {
			return true
		}
	}
	return false
}

func (instance *PathsHelper) IsDir(path string) (bool, error) {
	fi, err := os.Lstat(instance.Absolute(path))
	if nil != err {
		return false, err
	}
	return fi.Mode().IsDir(), nil
}

func (instance *PathsHelper) IsFile(path string) (bool, error) {
	fi, err := os.Lstat(instance.Absolute(path))
	if nil == err {
		return fi.Mode().IsRegular(), err
	} else {
		// path or file does not exists
		// just check if has extension
		if len(filepath.Ext(path)) > 0 {
			return true, nil
		}
		if strings.HasSuffix(path, _pathSeparator) {
			// is a directory
			return false, err
		}
		return true, err
	}
}

func (instance *PathsHelper) IsFilePath(path string) bool {
	if b, _ := instance.Exists(path); b {
		b, _ = instance.IsFile(path)
		return b
	} else {
		if strings.HasSuffix(path, _pathSeparator) {
			// is a directory
			return false
		}
		ext := filepath.Ext(path)
		return len(ext) > 0 && !strings.Contains(ext, " ")
	}
}

func (instance *PathsHelper) IsDirPath(path string) bool {
	if b, _ := instance.Exists(path); b {
		b, _ = instance.IsDir(path)
		return b
	} else {
		return len(filepath.Ext(path)) == 0
	}
}

func (instance *PathsHelper) IsAbs(path string) bool {
	if instance.IsUrl(path) {
		return true
	}
	return filepath.IsAbs(path)
}

func (instance *PathsHelper) IsUrl(path string) bool {
	return strings.Index(path, "http") == 0
}

func (instance *PathsHelper) IsSymLink(path string) (bool, error) {
	fi, err := os.Lstat(instance.Absolute(path))
	return fi.Mode()&os.ModeSymlink != 0, err
}

func (instance *PathsHelper) IsHiddenFile(path string) (bool, error) {
	return isHiddenFile(path)
}

func (instance *PathsHelper) IsSameFile(path1, path2 string) (bool, error) {
	f1, err1 := os.Lstat(instance.Absolute(path1))
	f2, err2 := os.Lstat(instance.Absolute(path2))
	if nil != err1 {
		return false, err1
	}
	if nil != err2 {
		return false, err2
	}
	return sameFile(f1, f2), nil
}

func (instance *PathsHelper) IsSameFileInfo(f1, f2 os.FileInfo) bool {
	return sameFile(f1, f2)
}

func (instance *PathsHelper) ListAll(root string) ([]string, error) {
	var response []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		response = append(response, path)
		return nil
	})
	return response, err
}

func (instance *PathsHelper) Walk(root, filter string, callback func(path string, info os.FileInfo) error) (err error) {
	if nil == callback {
		return
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if nil == err {
			if nil != info && !info.IsDir() {
				if len(filter) == 0 {
					return callback(path, info)
				} else {
					name := filepath.Base(path)
					if len(Regex.WildcardMatch(name, filter)) > 0 {
						return callback(path, info)
					}
				}
			} else {
				return callback(path, info)
			}
		}
		return err
	})
	return
}

func (instance *PathsHelper) PatternMatchBase(filename, pattern string) bool {
	name := filepath.Base(filename)
	return len(Regex.WildcardMatch(name, pattern)) > 0
}

func (instance *PathsHelper) ListFiles(root string, rawFilter interface{}) ([]string, error) {
	var response []string
	filters := make([]string, 0)
	if filter, ok := rawFilter.(string); ok {
		filters = append(filters, filter)
	} else if a, ok := rawFilter.([]string); ok {
		filters = append(filters, a...)
	}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if nil != info && !info.IsDir() {
			if len(filters) == 0 {
				response = append(response, path)
			} else {
				name := filepath.Base(path)
				for _, filter := range filters {
					if len(Regex.WildcardMatch(name, filter)) > 0 {
						response = append(response, path)
						break // added
					}
				}
			}
		}
		return nil
	})
	return response, err
}

func (instance *PathsHelper) WalkFiles(root string, filter string, callback func(path string) error) (err error) {
	if nil == callback {
		return
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if nil != info && !info.IsDir() {
			if len(filter) == 0 {
				return callback(path)
			} else {
				name := filepath.Base(path)
				if len(Regex.WildcardMatch(name, filter)) > 0 {
					return callback(path)
				}
			}
		}
		return nil
	})
	return
}

func (instance *PathsHelper) ListDir(root string) ([]string, error) {
	var response []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != root {
			response = append(response, path)
		}
		return nil
	})
	return response, err
}

func (instance *PathsHelper) WalkDir(root string, callback func(path string) error) (err error) {
	if nil == callback {
		return
	}
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != root {
			return callback(path)
		}
		return nil
	})
	return
}

func (instance *PathsHelper) ReadDirOnly(root string) ([]string, error) {
	var response []string
	infos, err := os.ReadDir(root)
	for _, info := range infos {
		if info.IsDir() {
			response = append(response, info.Name())
		}
	}
	return response, err
}

func (instance *PathsHelper) ReadFileOnly(root string) ([]string, error) {
	var response []string
	infos, err := os.ReadDir(root)
	for _, info := range infos {
		if !info.IsDir() {
			response = append(response, info.Name())
		}
	}
	return response, err
}

// ReadDir reads the directory named by dirname and returns
// a list of file's and directory's name.
func (instance *PathsHelper) ReadDir(root string) ([]string, error) {
	var response []string
	infos, err := os.ReadDir(root)
	for _, info := range infos {
		response = append(response, info.Name())
	}
	return response, err
}

func (instance *PathsHelper) WalkFilesOnOutput(root string, filter string, output io.Writer) (err error) {
	err = instance.WalkFiles(root, filter, func(path string) error {
		_, e := output.Write([]byte((path + "\n")))
		if nil != e {
			return e
		}
		return nil
	})
	return
}

func (instance *PathsHelper) WalkFilesOnOutputFile(root string, filter string, filename string) (err error) {
	var f *os.File
	if b, _ := Paths.Exists(filename); b {
		f, err = os.OpenFile(filename,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		f, err = os.Create(filename)
	}

	if nil == err {
		defer f.Close()
		w := bufio.NewWriter(f)
		err = instance.WalkFilesOnOutput(root, filter, w)
		_ = w.Flush()
	}
	return
}

func (instance *PathsHelper) TmpFileName(extension string) string {
	uuid := qb_rnd.Rnd.Uuid()
	if len(uuid) == 0 {
		uuid = "temp_file"
	}

	return uuid + ensureDot(extension)
}

func (instance *PathsHelper) TmpFile(extension string) string {
	path := filepath.Join(_temp_root, instance.TmpFileName(extension))
	return instance.Absolute(path)
}

func (instance *PathsHelper) ChangeFileName(fromPath, toFileName string) string {
	parent := filepath.Dir(fromPath)
	base := filepath.Base(fromPath)
	ext := filepath.Ext(base)
	if len(filepath.Ext(toFileName)) > 0 {
		return filepath.Join(parent, toFileName)
	}
	return filepath.Join(parent, toFileName+ext)
}

func (instance *PathsHelper) ChangeFileNameExtension(fromPath, toFileExtension string) string {
	parent := filepath.Dir(fromPath)
	base := filepath.Base(fromPath)
	ext := filepath.Ext(base)
	name := strings.Replace(base, ext, "", 1)
	return filepath.Join(parent, name+ensureDot(toFileExtension))
}

func (instance *PathsHelper) ChangeFileNameWithSuffix(path, suffix string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.Replace(base, ext, "", 1)
	return filepath.Join(filepath.Dir(path), name+suffix+ext)
}

func (instance *PathsHelper) ChangeFileNameWithPrefix(path, prefix string) string {
	base := filepath.Base(path)
	return filepath.Join(filepath.Dir(path), prefix+base)
}

func (instance *PathsHelper) NormalizePathForOS(path string) string {
	if qb_sys.Sys.IsWindows() {
		path = strings.ReplaceAll(path, "/", "\\")
	} else {
		path = strings.ReplaceAll(path, "\\", "/")
	}
	return instance.CleanPath(path)
}

func (instance *PathsHelper) Clean(pathOrUrl string) string {
	if instance.IsUrl(pathOrUrl) {
		return instance.CleanUrl(pathOrUrl)
	}
	return instance.CleanPath(pathOrUrl)
}

// CleanPath returns the shortest path name equivalent to path
// by purely lexical processing. It applies the following rules
// iteratively until no further processing can be done:
//
//  1. Replace multiple Separator elements with a single one.
//  2. Eliminate each . path name element (the current directory).
//  3. Eliminate each inner .. path name element (the parent directory)
//     along with the non-.. element that precedes it.
//  4. Eliminate .. elements that begin a rooted path:
//     that is, replace "/.." by "/" at the beginning of a path,
//     assuming Separator is '/'.
//
// The returned path ends in a slash only if it represents a root directory,
// such as "/" on Unix or `C:\` on Windows.
//
// Finally, any occurrences of slash are replaced by Separator.
//
// If the result of this process is an empty string, Clean
// returns the string ".".
//
// See also Rob Pike, “Lexical File Names in Plan 9 or
// Getting Dot-Dot Right,”
// https://9p.io/sys/doc/lexnames.html
func (instance *PathsHelper) CleanPath(p string) string {
	return filepath.Clean(p)
}

// CleanUrl is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//  1. Replace multiple slashes with a single slash.
//  2. Eliminate each . path name element (the current directory).
//  3. Eliminate each inner .. path name element (the parent directory)
//     along with the non-.. element that precedes it.
//  4. Eliminate .. elements that begin a rooted path:
//     that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned
func (instance *PathsHelper) CleanUrl(p string) string {
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	trailing := n > 2 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)

	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r++

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 2

			if w > 1 {
				// can backtrack
				w--

				if buf == nil {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// real path element.
			// add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	if buf == nil {
		return p[:w]
	}
	return string(buf[:w])
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func ensureDot(extension string) string {
	if strings.Index(extension, ".") == -1 {
		extension = "." + extension
	}
	return extension
}

// internal helper to lazily create a buffer if necessary
func bufApp(buf *[]byte, s string, w int, c byte) {
	if *buf == nil {
		if s[w] == c {
			return
		}

		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	}
	(*buf)[w] = c
}
