# Scheduler #

Schedule everything.

## Configuration ##

Scheduler is fully configurable.

```
{
    "uid": "Sample Scheduler",
    "sync": false,
    "schedules": [
      {
        "uid": "sample_1",
        "start_at": "",
        "timeline": "second:30"
      },
      {
        "uid": "sample_2",
        "start_at": "11:00:00",
        "timeline": "hour:12"
      }
    ]
}
```

**Parameters**

* uid: A name for your scheduler (used in logs).
* sync: Default is False. Enable sync mode if you need that OnSchedule events are locking.
* schedules: Array of `schedule` objects. You can have more than one single scheduled job. 
    * uid: A name used to identify a schedule in the array
    * start_at: (Optional) Time you want to launch e trigger. 
    * timeline: Duration of your scheduled cycle. Do you want to launch a trigger every 5 minutes? Just set the timeline to "minute:5" 



## Sample Code ##
```
    sched := lygo_scheduler.NewSchedulerFromFile("./scheduler.json")
	if sched.HasErrors(){
		fmt.Println("Parsing Errors", sched.GetErrors())
	} else {
        sched.OnSchedule(func(schedule *lygo_scheduler.SchedulerTask) {
        	// do something useful here....
            fmt.Println("EVENT", schedule.Settings().Uid, schedule.Settings().StartAt, schedule.Settings().Timeline)
        })
        sched.Start()
    }
```