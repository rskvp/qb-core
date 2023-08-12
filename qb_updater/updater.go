package qb_updater

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rskvp/qb-core/qb_events"
	"github.com/rskvp/qb-core/qb_rnd"
	"github.com/rskvp/qb-core/qb_scheduler"
	"github.com/rskvp/qb-core/qb_utils"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

var (
	ErrorMissingConfigurationParameter = errors.New("missing_configuration_parameter_error")
)

const (
	VariableDirHome  = "$dir_home"  // root
	VariableDirStart = "$dir_start" // binary launch dir
	VariableDirApp   = "$dir_app"   // binary dir
	VariableDirWork  = "$dir_work"  // workspace

	onUpgrade       = "on_upgrade"
	onError         = "on_error"
	onTask          = "on_task"
	onRelaunch      = "on_relaunch"
	onStartLauncher = "on_start_launcher"
	onStopLauncher  = "on_stop_laucher"

	DirStart = "start"
	DirApp   = "app"
	DirWork  = "*"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type GenericEventHandler func(updater *Updater, eventName string, args []interface{})
type UpdaterErrorHandler func(err string)
type UpdaterUpgradeHandler func(fromVersion, toVersion string, files []string)
type LauncherStartHandler func(command string)
type LauncherStartedHandler func(command string, pid int)
type LauncherQuitHandler func(command string, pid int)
type TaskHandler func(taskUID string, payload map[string]interface{})

type Updater struct {
	root                  string
	dirStart              string
	dirApp                string
	dirWork               string
	uid                   string
	settings              *Settings
	variables             map[string]string
	launcher              *Launcher
	schedulerUpdate       *qb_scheduler.Scheduler
	schedulerRestart      *qb_scheduler.Scheduler
	schedulerTask         *qb_scheduler.Scheduler
	events                *qb_events.Emitter
	genericHandlers       []GenericEventHandler
	errorHandlers         []UpdaterErrorHandler
	upgradeHandlers       []UpdaterUpgradeHandler
	launchStartHandlers   []LauncherStartHandler
	launchStartedHandlers []LauncherStartedHandler
	launchQuitHandlers    []LauncherQuitHandler
	taskHandlers          []TaskHandler
	chanQuit              chan bool
	started               bool
	isReadyToRestart      bool // is launcher already started al least once?
	processMux            sync.Mutex

	// state
	_isUpdating               bool
	_launcherStoppedForUpdate bool
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewUpdater(settings ...interface{}) *Updater {
	instance := new(Updater)
	instance.uid = qb_rnd.Rnd.Uuid()
	instance.started = false
	instance.root = qb_utils.Paths.Absolute("./")
	instance.dirStart = qb_utils.Paths.GetWorkspace(DirStart).GetPath()
	instance.dirApp = qb_utils.Paths.GetWorkspace(DirApp).GetPath()
	instance.dirWork = qb_utils.Paths.GetWorkspace(DirWork).GetPath()
	instance.chanQuit = make(chan bool, 1)
	instance.events = qb_events.Events.NewEmitter()
	instance.variables = make(map[string]string)

	if len(settings) > 0 {
		instance.init(settings[0])
	} else {
		instance.init("./updater.json")
	}

	if nil == instance.settings {
		instance.settings = new(Settings)
		instance.settings.ScheduledUpdates = make([]*qb_scheduler.Schedule, 0)
	}

	instance.genericHandlers = make([]GenericEventHandler, 0)
	instance.errorHandlers = make([]UpdaterErrorHandler, 0)
	instance.upgradeHandlers = make([]UpdaterUpgradeHandler, 0)
	instance.initUpdaterEvents()

	instance.launchStartHandlers = make([]LauncherStartHandler, 0)
	instance.launchStartedHandlers = make([]LauncherStartedHandler, 0)
	instance.launchQuitHandlers = make([]LauncherQuitHandler, 0)
	instance.taskHandlers = make([]TaskHandler, 0)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Updater) Settings() *Settings {
	if nil != instance {
		return instance.settings
	}
	return nil
}

func (instance *Updater) SetUid(uid string) {
	if nil != instance {
		instance.uid = uid
	}
}

func (instance *Updater) GetUid() string {
	if nil != instance {
		return instance.uid
	}
	return ""
}

func (instance *Updater) SetRoot(path string) {
	if nil != instance {
		instance.root = qb_utils.Paths.Absolute(path)
	}
}

func (instance *Updater) GetRoot() string {
	if nil != instance {
		return instance.root
	}
	return ""
}

func (instance *Updater) GetDirStart() string {
	if nil != instance {
		return instance.dirStart
	}
	return ""
}

func (instance *Updater) GetDirApp() string {
	if nil != instance {
		return instance.dirApp
	}
	return ""
}

func (instance *Updater) GetDirWork() string {
	if nil != instance {
		return instance.dirWork
	}
	return ""
}

func (instance *Updater) GetSettingVersionFile() string {
	if nil != instance {
		if nil != instance.settings {
			return instance.settings.VersionFile
		}
	}
	return ""
}

func (instance *Updater) GetSettingKeepAlive() bool {
	if nil != instance {
		if nil != instance.settings {
			return instance.settings.KeepAlive
		}
	}
	return false
}

func (instance *Updater) SetVariable(name, value string) {
	if nil != instance && nil != instance.variables {
		if strings.Index(name, "$") == -1 {
			name = "$" + name
		}
		instance.variables[name] = value
	}
}

func (instance *Updater) GetVariable(name string) string {
	if nil != instance && nil != instance.variables {
		if strings.Index(name, "$") == -1 {
			name = "$" + name
		}
		if v, b := instance.variables[name]; b {
			return v
		}
	}
	return ""
}

func (instance *Updater) OnEvent(handler GenericEventHandler) {
	if nil != instance && nil != handler {
		instance.genericHandlers = append(instance.genericHandlers, handler)
	}
}

func (instance *Updater) OnError(handler UpdaterErrorHandler) {
	if nil != instance && nil != handler {
		instance.errorHandlers = append(instance.errorHandlers, handler)
	}
}

func (instance *Updater) OnUpgrade(handler UpdaterUpgradeHandler) {
	if nil != instance && nil != handler {
		instance.upgradeHandlers = append(instance.upgradeHandlers, handler)
	}
}

func (instance *Updater) OnTask(handler TaskHandler) {
	if nil != instance && nil != handler {
		instance.taskHandlers = append(instance.taskHandlers, handler)
	}
}

func (instance *Updater) OnLaunchStart(handler LauncherStartHandler) {
	if nil != instance && nil != handler {
		instance.launchStartHandlers = append(instance.launchStartHandlers, handler)
	}
}

func (instance *Updater) OnLaunchStarted(handler LauncherStartedHandler) {
	if nil != instance && nil != handler {
		instance.launchStartedHandlers = append(instance.launchStartedHandlers, handler)
	}
}

func (instance *Updater) OnLaunchQuit(handler LauncherQuitHandler) {
	if nil != instance && nil != handler {
		instance.launchQuitHandlers = append(instance.launchQuitHandlers, handler)
	}
}

func (instance *Updater) HasUpdates() bool {
	currentVersion, remoteVersion, _ := instance.getVersions()
	return instance.needUpdate(currentVersion, remoteVersion)
}

func (instance *Updater) IsUpgradable(currentVersion, remoteVersion string) bool {
	return instance.needUpdate(currentVersion, remoteVersion)
}

func (instance *Updater) IsProcessRunning() bool {
	if nil != instance && nil != instance.launcher {
		instance.processMux.Lock()
		defer instance.processMux.Unlock()

		return instance.launcher.Pid() > -1
	}
	return false
}

func (instance *Updater) GetProcessOutput() string {
	if nil != instance && instance.IsProcessRunning() {
		return instance.launcher.Output()
	}
	return ""
}

func (instance *Updater) GetProcessPid() int {
	if nil != instance && nil != instance.launcher {
		instance.processMux.Lock()
		defer instance.processMux.Unlock()

		return instance.launcher.Pid()
	}
	return -1
}

func (instance *Updater) IsUpdating() bool {
	if nil != instance {
		return instance._isUpdating
	}
	return false
}

// Start
// run the update job and also check immediately for updates
func (instance *Updater) Start() (updated bool, fromVersion string, toVersion string, files []string, err error) {
	if nil != instance {
		instance.processMux.Lock()
		defer instance.processMux.Unlock()

		instance.started = true
		instance.refreshVariables()

		instance._isUpdating = true // BEGIN UPDATING STATE

		// is launcher active?
		if nil != instance.launcher {
			if instance.HasUpdates() {
				instance.stopLauncher()
			}
		}

		// check updates and download
		if len(instance.GetSettingVersionFile()) > 0 {
			updated, fromVersion, toVersion, files, err = instance.checkUpdates()
		}

		instance._isUpdating = false // END UPDATING STATE

		if len(instance.settings.CommandToRun) > 0 {
			// LAUNCH PROGRAM
			launchErr := instance.startLauncher()
			if nil != launchErr {
				instance.events.EmitAsync(onError, launchErr.Error())
				if nil == err {
					err = launchErr
				}
			}
		}

		if nil != err {
			instance.events.EmitAsync(onError, err.Error())
		}
		if updated {
			instance.events.EmitAsync(onUpgrade, fromVersion, toVersion, files)
		}

		// START SCHEDULER IF ANY AND IF NOT STARTED YET
		// schedulerUpdate should always start because update server may be down,
		// but updater should continue to check for new updates
		instance.initScheduler()
	}
	return updated, fromVersion, toVersion, files, err
}

func (instance *Updater) Stop() {
	if nil != instance {
		instance.processMux.Lock()
		defer instance.processMux.Unlock()

		if nil != instance.launcher {
			_ = instance.launcher.Kill() // STOP RUNNING PROGRAM
		}
		instance.started = false
		instance.chanQuit <- true
	}
}

func (instance *Updater) Wait() {
	if nil != instance {
		if !instance.started {
			_, _, _, _, _ = instance.Start()
		}
		<-instance.chanQuit // wait exit
		instance.chanQuit = make(chan bool, 1)
	}
}

func (instance *Updater) ReStart() {
	if nil != instance && nil != instance.launcher {
		instance.processMux.Lock()
		defer instance.processMux.Unlock()

		if !instance.isReadyToRestart {
			instance.isReadyToRestart = true
		} else {
			// temporary stop for updater scheduler
			instance.schedulerUpdate.Pause()
			defer instance.schedulerUpdate.Resume()

			// stop the launcher
			instance.stopLauncher()

			// start launcher again
			launchErr := instance.startLauncher()
			if nil != launchErr {
				instance.events.EmitAsync(onError, launchErr.Error())
			}
		}
	}
}

func (instance *Updater) ReLaunch() {
	if nil != instance && nil != instance.launcher {
		instance.processMux.Lock()
		defer instance.processMux.Unlock()

		// temporary stop for updater scheduler
		instance.schedulerUpdate.Pause()
		defer instance.schedulerUpdate.Resume()

		cause := instance.getLauncher().Error()

		// stop the launcher
		instance.stopLauncher()

		// start launcher again
		launchErr := instance.startLauncher()
		if nil != launchErr {
			instance.events.EmitAsync(onError, launchErr.Error())
		}
		instance.events.EmitAsync(onRelaunch, cause)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Updater) init(settings interface{}) {
	if nil != settings {
		if v, b := settings.(string); b {
			if qb_utils.Regex.IsValidJsonObject(v) {
				// json string
				var s Settings
				err := qb_utils.JSON.Read(v, &s)
				if nil == err {
					instance.settings = &s
				}
			} else if b, _ := qb_utils.Paths.Exists(v); b {
				// load from file
				text, err := qb_utils.IO.ReadTextFromFile(v)
				if nil == err {
					instance.init(text)
				}
			}
		} else if v, b := settings.(Settings); b {
			instance.settings = &v
		}
	}
}

func (instance *Updater) initUpdaterEvents() {
	instance.events.On(onError, func(event *qb_events.Event) {
		if nil != instance && nil != instance.errorHandlers {
			var err string
			item := event.Argument(0)
			if v, b := item.(string); b {
				err = v
			} else if v, b := item.(error); b {
				err = v.Error()
			}
			if len(err) > 0 {
				for _, handler := range instance.errorHandlers {
					if nil != handler {
						handler(err)
					}
				}
				instance.bubbleGenericEvent(onError, err)
			}
		}
	})
	instance.events.On(onUpgrade, func(event *qb_events.Event) {
		if nil != instance && nil != instance.upgradeHandlers {
			arg1 := event.Argument(0)
			arg2 := event.Argument(1)
			arg3 := event.Argument(2)
			if fromVersion, b := arg1.(string); b {
				if toVersion, b := arg2.(string); b {
					files := qb_utils.Convert.ToArrayOfString(arg3)
					for _, handler := range instance.upgradeHandlers {
						if nil != handler {
							handler(fromVersion, toVersion, files)
						}
					}

					instance.bubbleGenericEvent(onUpgrade, fromVersion, toVersion, files)
				}
			}
		}
	})
	instance.events.On(onTask, func(event *qb_events.Event) {
		if nil != instance && nil != instance.taskHandlers {
			arg1 := event.Argument(0)
			if task, b := arg1.(*qb_scheduler.SchedulerTask); b {
				uid := task.Uid
				payload := task.Payload
				for _, handler := range instance.taskHandlers {
					if nil != handler {
						handler(uid, payload)
					}
				}
				instance.bubbleGenericEvent(onTask, uid, payload)
			}
		}
	})
}

func (instance *Updater) startLauncher() (err error) {
	if nil != instance {
		// LAUNCH PROGRAM
		launcher := instance.getLauncher()
		if launcher.Pid() == -1 {
			// service was closed: run again
			err = instance.getLauncher().Run(
				replaceVars(instance.settings.CommandToRun, instance.variables),
			)

			instance._launcherStoppedForUpdate = false
		}
	}
	return
}

func (instance *Updater) stopLauncher() {
	if nil != instance && nil != instance.launcher {
		instance._launcherStoppedForUpdate = true

		_ = instance.launcher.Kill() // STOP RUNNING PROGRAM
		instance.launcher = nil
	}
}

func (instance *Updater) getLauncher() *Launcher {
	if nil == instance.launcher {
		// the launcher always keep the program session
		instance.launcher = NewLauncher(true)
		instance.initLauncherEvents()
	}
	return instance.launcher
}

func (instance *Updater) initLauncherEvents() {
	if nil != instance.launcher {
		instance.launcher.OnStart(func(command string) {
			if nil != instance && nil != instance.launchStartHandlers {
				for _, callback := range instance.launchStartHandlers {
					callback(command)
				}
				instance.bubbleGenericEvent(onStart, command)
			}
		})
		instance.launcher.OnStarted(func(command string, pid int) {
			if nil != instance && nil != instance.launchStartedHandlers {
				for _, callback := range instance.launchStartedHandlers {
					callback(command, pid)
				}
				instance.bubbleGenericEvent(onStarted, command, pid)
			}
		})
		instance.launcher.OnQuit(func(command string, pid int) {
			if nil != instance && nil != instance.launchQuitHandlers {
				for _, callback := range instance.launchQuitHandlers {
					callback(command, pid)
				}
				instance.bubbleGenericEvent(onQuit, command, pid)
			}
			instance.checkIfKeepAlive()
		})
	}
}

func (instance *Updater) bubbleGenericEvent(eventName string, args ...interface{}) {
	for _, handler := range instance.genericHandlers {
		if nil != handler {
			handler(instance, eventName, args)
		}
	}
}

func (instance *Updater) initScheduler() {
	if nil != instance {
		// UPDATE
		instance.initSchedulerUpdate()

		// RESTART
		instance.initSchedulerRestart()

		// TASKS
		instance.initSchedulerTasks()
	}
}

func (instance *Updater) initSchedulerUpdate() {
	if nil != instance && nil == instance.schedulerUpdate {
		// schedulerUpdate not already initialized
		if len(instance.settings.ScheduledUpdates) > 0 {
			instance.schedulerUpdate = qb_scheduler.NewScheduler()
			for _, schedule := range instance.settings.ScheduledUpdates {
				instance.schedulerUpdate.AddSchedule(schedule)
			}
			instance.schedulerUpdate.OnError(func(error string) {
				instance.events.EmitAsync(onError, error)
			})
			instance.schedulerUpdate.OnSchedule(func(schedule *qb_scheduler.SchedulerTask) {
				// fmt.Println("instance.schedulerUpdate.OnSchedule", schedule.String())
				_, _, _, _, _ = instance.Start() // check updates
			})
			instance.schedulerUpdate.Start()
		}
	}
}

func (instance *Updater) initSchedulerRestart() {
	if nil != instance && nil == instance.schedulerRestart {
		// schedulerUpdate not already initialized
		if len(instance.settings.ScheduledRestart) > 0 {
			instance.schedulerRestart = qb_scheduler.NewScheduler()
			for _, schedule := range instance.settings.ScheduledRestart {
				instance.schedulerRestart.AddSchedule(schedule)
			}
			instance.schedulerRestart.OnError(func(error string) {
				instance.events.EmitAsync(onError, error)
			})
			instance.schedulerRestart.OnSchedule(func(schedule *qb_scheduler.SchedulerTask) {
				instance.ReStart() // restart application
			})
			instance.schedulerRestart.Start()
		}
	}
}

func (instance *Updater) initSchedulerTasks() {
	if nil != instance && nil == instance.schedulerTask {
		// schedulerUpdate not already initialized
		if len(instance.settings.ScheduledTasks) > 0 {
			instance.schedulerTask = qb_scheduler.NewScheduler()
			for _, schedule := range instance.settings.ScheduledTasks {
				instance.schedulerTask.AddSchedule(schedule)
			}
			instance.schedulerTask.OnError(func(error string) {
				instance.events.EmitAsync(onError, error)
			})
			instance.schedulerTask.OnSchedule(func(schedule *qb_scheduler.SchedulerTask) {
				instance.events.EmitAsync(onTask, schedule)
			})
			instance.schedulerTask.Start()
		}
	}
}

func (instance *Updater) refreshVariables() {
	if nil != instance {
		instance.variables[VariableDirHome] = instance.root
		instance.variables[VariableDirStart] = instance.dirStart
		instance.variables[VariableDirApp] = instance.dirApp
		instance.variables[VariableDirWork] = instance.dirWork
	}
}

func (instance *Updater) checkUpdates() (updated bool, currentVersion string, remoteVersion string, files []string, err error) {
	if len(instance.settings.VersionFile) > 0 {
		updated, currentVersion, remoteVersion, files, err = instance.check()
	} else {
		err = qb_utils.Errors.Prefix(ErrorMissingConfigurationParameter, "Missing Configuration Parameter 'VersionFile': ")
	}
	return updated, currentVersion, remoteVersion, files, err
}

func (instance *Updater) getVersions() (currentVersion string, remoteVersion string, filename string) {
	if nil != instance.settings && len(instance.settings.VersionFile) > 0 {
		url := instance.settings.VersionFile
		filename = qb_utils.Paths.Concat(instance.root, qb_utils.Paths.FileName(url, true))

		currentVersion = instance.getCurrentVersion(filename)
		remoteVersion = instance.getRemoteVersion(url)
	}
	return currentVersion, remoteVersion, filename
}

func (instance *Updater) check() (bool, string, string, []string, error) {
	currentVersion, remoteVersion, filename := instance.getVersions()
	needUpdate := instance.needUpdate(currentVersion, remoteVersion)
	files := make([]string, 0)

	if needUpdate {
		// download & install packages
		for _, v := range instance.settings.PackageFiles {
			source := v.File // may be an URL too.
			target := qb_utils.Paths.Absolute(replaceVars(v.Target, instance.variables))
			err := instance.install(source, target)
			if nil != err {
				return false, currentVersion, remoteVersion, files, err
			} else {
				files = append(files, source, target)
			}
		}
	}

	if len(currentVersion) == 0 || needUpdate {
		// update version file
		_, err := qb_utils.IO.WriteTextToFile(remoteVersion, filename)
		if nil != err {
			return false, currentVersion, remoteVersion, files, err
		}
	}

	return needUpdate, currentVersion, remoteVersion, files, nil
}

func (instance *Updater) getCurrentVersion(filename string) string {
	if s, err := qb_utils.IO.ReadTextFromFile(filename); nil == err {
		return strings.Trim(s, " \n")
	}
	// file does not exists
	if instance.settings.VersionFileRequired {
		return "0"
	}
	return ""
}

func (instance *Updater) getRemoteVersion(filename string) string {
	data, err := instance.download(filename)
	if nil == err {
		return strings.Trim(string(data), " \n")
	}
	return ""
}

func (instance *Updater) install(url, target string) error {
	data, err := instance.download(url)
	if nil != err {
		return err
	}

	filename := qb_utils.Paths.FileName(url, true)
	ext := qb_utils.Paths.ExtensionName(filename)
	targetDir := qb_utils.Paths.Absolute(target) + string(os.PathSeparator)
	err = qb_utils.Paths.Mkdir(targetDir) // creates sub folders if missing
	if nil != err {
		return err
	}

	if "zip" == ext {
		// unzip
		tmp := qb_rnd.Rnd.Uuid() + ".zip"
		defer qb_utils.IO.RemoveSilent(tmp)
		_, err := qb_utils.IO.WriteBytesToFile(data, tmp)
		if nil != err {
			return err
		}
		_, err = qb_utils.Zip.Unzip(tmp, targetDir)
		if nil != err {
			return err
		}
	} else {
		// copy
		_, err := qb_utils.IO.WriteBytesToFile(data, qb_utils.Paths.Concat(targetDir, filename))
		if nil != err {
			return err
		}
	}

	return nil
}

func (instance *Updater) download(url string) ([]byte, error) {
	if len(url) > 0 {
		if strings.Index(url, "http") > -1 {
			// HTTP
			tr := &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    15 * time.Second,
				DisableCompression: true,
			}
			client := &http.Client{Transport: tr}
			resp, err := client.Get(url)
			if nil == err {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if nil == err {
					return body, nil
				} else {
					return []byte{}, err
				}
			} else {
				return []byte{}, err
			}
		} else {
			// FILE SYSTEM
			path := replaceVars(url, instance.variables)
			return qb_utils.IO.ReadBytesFromFile(path)
		}
	}
	return []byte{}, qb_utils.Errors.Prefix(ErrorMissingConfigurationParameter, "Missing Configuration Parameter 'VersionFile': ")
}

func (instance *Updater) needUpdate(currentVersion, remoteVersion string) bool {
	if len(currentVersion) > 0 && len(remoteVersion) > 0 && currentVersion != remoteVersion {
		v1 := strings.Split(strings.Trim(currentVersion, " \n"), ".")
		v2 := strings.Split(strings.Trim(remoteVersion, " \n"), ".")
		for i := 0; i < 3; i++ {
			a := qb_utils.Convert.ToInt(qb_utils.Arrays.GetAt(v1, i, 0))
			b := qb_utils.Convert.ToInt(qb_utils.Arrays.GetAt(v2, i, 0))
			if b == a {
				continue
			} else if b < a {
				return false
			} else if b > a {
				return true
			}
		}
	}
	return false
}

func (instance *Updater) checkIfKeepAlive() {
	if nil != instance && instance.GetSettingKeepAlive() && !instance.IsUpdating() {
		instance.ReLaunch()
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func replaceVars(text string, variables map[string]string) string {
	txt := text
	for k, v := range variables {
		txt = strings.Replace(txt, k, v, -1)
	}
	return txt
}
