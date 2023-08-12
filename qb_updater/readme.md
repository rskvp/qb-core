# Updater

Update and launch a program.
Updater is a simple "live-update" utility that can update (searching on web for latest version) and launch a program.

When Updater starts, first read the latest version from remote "version.txt" file.
If remote version is newer than the local version, Updater start download and install (unzip) new version.
Finally, launch the program.

Optionally you can schedule remote check and auto-updates.

## Start and Update

Updater standard behaviour is "update-and-start".

Every time you start a program using Updater, a version check is performed before launch your program.

## Scheduled Checks

Scheduled checks keep your application updated.

At every update the program is stopped and restarted. Do not schedule updates if you do not want stop your program.

For scheduling tasks Updater uses: [Scheduler](../qb_updater/readme.md)

## Configuration

Updater is fully configurable changing file updater.json.

```
{
  "version_file": "https://gianangelogeminiani.me/download/version.txt",
  "package_files": [
    {
      "file": "https://gianangelogeminiani.me/download/package.zip",
      "target": "./bin"
    }
  ],
  "command_to_run": "./bin/myservice",
  "scheduled_updates": [
    {
      "uid": "every_3_seconds",
      "start_at": "",
      "timeline": "second:3"
    }
  ]
}
```

**Parameters**

- version_file: Path (relative or absolute) to text file containing latest version number.
- package_files: Array of objects (PackageFile) to download and unzip (if archive). PackageFile contains "file" and "target" fields.
- command_to_run: Command to run when screen launcher is active. Use this to run your program.

**Variables**

Some parameters (`command_to_run`, `package_files.target`) can contain variables.

- $dir_home: Is replaced with Application absolute path.

## Sample Code

```
    updater := qb_updater.NewUpdater()
	updater.Settings().VersionFile = "./versions/version.txt"
	updater.Settings().PackageFiles = make([]*lygo_updater.PackageFile, 0)
	updater.Settings().PackageFiles = append(updater.Settings().PackageFiles, &lygo_updater.PackageFile{
		File:   "./versions/test_fiber.zip",
		Target: "./versions_install",
	})
	updater.Settings().CommandToRun = "./versions_install/test_fiber"
	updater.OnError(func(err string) {
		fmt.Println("ERROR", err)
	})
	count:=0
	updater.OnUpgrade(func(fromVersion, toVersion string) {
		count++
		now := time.Now()
		t := fmt.Sprintf("%v:%v:%v", now.Hour(), now.Minute(), now.Second())
		fmt.Println("UPDATE", count, t, "\t" + fromVersion + " -> " + toVersion, "pid:", updater.GetProcessPid())
	})

    // use Start method to check for updates and go on
	// updater.Start()

    // use Wait method to Start and wait
    updater.Wait()

```
