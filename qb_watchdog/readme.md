# Watchdog file watcher

FileWatcher is a fork of [watcher](https://github.com/radovskyb/watcher) excellent project of Benjamin Radovsky.

`FileWatcher` is a Go package for watching for files or directory changes (recursively or non recursively) without using filesystem events, which allows it to work cross platform consistently.

`FileWatcher` watches for changes and notifies over channels either anytime an event or an error has occurred.

Events contain the `os.FileInfo` of the file or directory that the event is based on and the type of event and file or directory path.

[Features](#features)  
[Example](#example)  
[Contributing](#contributing)  
[Watcher Command](#command)

# Features

- Customizable polling interval.
- Filter Events.
- Watch folders recursively or non-recursively.
- Choose to ignore hidden files.
- Choose to ignore specified files and folders.
- Notifies the `os.FileInfo` of the file that the event is based on. e.g `Name`, `ModTime`, `IsDir`, etc.
- Notifies the full path of the file that the event is based on or the old and new paths if the event was a `Rename` or `Move` event.
- Limit amount of events that can be received per watching cycle.
- List the files being watched.
- Trigger custom events.

# Example

```go
package main

import (
	"fmt"
	"bitbucket.org/lygo/lygo_file_watcher"
"log"
	"regexp"
"time"


)

func main() {
	w := qb_watchdog.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)

	// Only notify rename and move events.
	w.FilterOps(lygo_file_watcher.Rename, lygo_file_watcher.Move)

	// Only files that match the regular expression during file listings
	// will be watched.
	r := regexp.MustCompile("^abc$")
	w.AddFilterHook(lygo_file_watcher.RegexFilterHook(r, false))

	go func() {
		for {
			select {
			case event := <-w.Event:
				fmt.Println(event) // Print the event's info.
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	// Watch this folder for changes.
	if err := w.Add("."); err != nil {
		log.Fatalln(err)
	}

	// Watch test_folder recursively for changes.
	if err := w.AddRecursive("../test_folder"); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	fmt.Println()

	// Trigger 2 events after watcher started.
	go func() {
		w.Wait()
		w.TriggerEvent(lygo_file_watcher.Create, nil)
		w.TriggerEvent(lygo_file_watcher.Remove, nil)
	}()

	// Start the watching process - it'll check for changes every 100ms.
	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}
```
