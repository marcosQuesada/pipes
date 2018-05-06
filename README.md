Golang Pipes
============

 Just a producer-consumer playground

 ## Application flow
  - Chain pipe:
    - Implements Producer - Consumer relation on a synchronous widget channel
    - Producer - Consumer Plugable Implementation
  - Producer: creates widgets and inject them on widget channel
  - Consumer: consume widgets from widget channel

 ## Concurrent Consumer pool
  - Implements Consumer interface, same as other consumers
  - builds a consumer pool creating consumers using Consumer Builder

 ## Threading model
  - consumers are delegate to its own goroutine
  - wait groups has been used to synchronize task execution

 ## Test suite:
  Pipe tests has been used under development and represent use cases
  You can run full test suite as:

```
 go test --race -v ./...
```

 ## Use Cases:

```
Single producer on single consumer
go run main.go -n 4
[widget_id_0 19:52:27.711509] consumer_0 131.967µs
[widget_id_1 19:52:27.711592] consumer_0 64.346µs
[widget_id_2 19:52:27.711662] consumer_0 2.89µs
[widget_id_3 19:52:27.711663] consumer_0 6.596µs

Single producer on delayed consumer
time go run main.go -n 4 -d 1000
[widget_id_0 19:52:34.706439] consumer_0 147.927µs
[widget_id_1 19:52:34.706574] consumer_0 1.000166327s
[widget_id_2 19:52:35.706758] consumer_0 1.000156113s
[widget_id_3 19:52:36.706937] consumer_0 1.000257372s
go run main.go -n 4 -d 1000  0,22s user 0,03s system 5% cpu 4,211 total

Single producer on concurrent consumer pool from delayed consumer
go run main.go -n 4 -d 1000 -c 3
[widget_id_0 19:52:41.487030] consumer_2 214.3µs
[widget_id_2 19:52:41.487143] consumer_1 124.674µs
[widget_id_1 19:52:41.487135] consumer_0 144.825µs
[widget_id_3 19:52:41.487144] consumer_2 1.000356027s

Full execution using blocking producers:
go run main.go  -n 10000 | wc -l
10000

Full execution using non blocking producers:
go run main.go  -x -n 10000 | wc -l
118

Execution on a discarded overrun buffer, widgets are stored on buffer, buffer is processed once non blocking insertion is done.
go run main.go  -b 4 -n 4 -d 1000
[widget_id_0 19:53:12.523533] consumer_0 203.38µs
[widget_id_1 19:53:12.523535] consumer_0 1.00044999s
[widget_id_2 19:53:12.523536] consumer_0 2.000685706s
[widget_id_3 19:53:12.523537] consumer_0 3.000889201s

```



