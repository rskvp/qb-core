package qb_updater

import "github.com/rskvp/qb-core/qb_scheduler"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type Settings struct {
	Uid                 string                   `json:"uid"`
	KeepAlive           bool                     `json:"keep_alive"`            // launch again if program is closed
	VersionFileRequired bool                     `json:"version_file_required"` // if true, first start will update all if version file does not exists
	VersionFile         string                   `json:"version_file"`
	PackageFiles        []*PackageFile           `json:"package_files"`
	CommandToRun        string                   `json:"command_to_run"`
	ScheduledUpdates    []*qb_scheduler.Schedule `json:"scheduled_updates"`
	ScheduledRestart    []*qb_scheduler.Schedule `json:"scheduled_restart"`
	ScheduledTasks      []*qb_scheduler.Schedule `json:"scheduled_tasks"`
}

type PackageFile struct {
	File   string `json:"file"`
	Target string `json:"target"`
}
