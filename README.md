docker-cleaner
-------
[![GitHub Workflow Status (branch)](https://img.shields.io/github/workflow/status/foxdalas/docker-cleaner/build-and-test/master?style=for-the-badge)](https://github.com/foxdalas/docker-cleaner/actions)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/foxdalas/docker-cleaner?style=for-the-badge)](https://github.com/foxdalas/docker-cleaner/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/foxdalas/docker-cleaner?style=for-the-badge)](https://hub.docker.com/repository/docker/foxdalas/docker-cleaner)

One more docker cleaner for docker daemon. This tool working like ```docker system prune -a``` with prometheus metrics.

Usage patterns
-------
* build agents
* docker hosts with many deploys

Flags
-------
```
Usage of ./docker-cleaner:
  -cleaner.interval duration
        Cleaner check interval (default 15s)
  -docker.dir string
        Docker storage directory (default "/var/lib/docker")
  -docker.threshold float
        Docker volume usage threshold (default 50)
  -docker.ttl duration
        Docker volumes TTL. Same until=48h (default 48h0m0s)
  -exporter.host string
        Docker cleaner exporter listen host (default "0.0.0.0")
  -exporter.port int
        Docker cleaner exporter listen port (default 9203)
  -exporter.telemetry-path string
        Docker cleaner exporter path under which to expose metrics. (default "/metrics")
  -exporter.timeout duration
        Docker cleaner exporter timeout (default 15s)
```

ENVS
-------
```
LOG_LEVEL = (debug|info|warn|error)
LOG_TYPE = (text|json)
```

Metrics
-------
```
# HELP docker_cleaner_disk_reclaimable Docker daemon disk reclaimable
# TYPE docker_cleaner_disk_reclaimable gauge
docker_cleaner_disk_reclaimable{type="build_cache"} 199
docker_cleaner_disk_reclaimable{type="containers"} 0
docker_cleaner_disk_reclaimable{type="images"} 2.0836973e+07
docker_cleaner_disk_reclaimable{type="volumes"} 0
# HELP docker_cleaner_disk_usage Docker daemon disk usage
# TYPE docker_cleaner_disk_usage gauge
docker_cleaner_disk_usage{type="build_cache"} 199
docker_cleaner_disk_usage{type="containers"} 0
docker_cleaner_disk_usage{type="images"} 2.0836973e+07
docker_cleaner_disk_usage{type="total"} 6.84964343808e+11
docker_cleaner_disk_usage{type="volumes"} 0
# HELP docker_cleaner_disk_usage_percents Docker daemon disk usage.
# TYPE docker_cleaner_disk_usage_percents gauge
docker_cleaner_disk_usage_percents 68.47993321066149
# HELP docker_cleaner_last Cleaner last run.
# TYPE docker_cleaner_last counter
docker_cleaner_last 0
# HELP docker_cleaner_up Could the docker-cleaner server be reached.
# TYPE docker_cleaner_up gauge
docker_cleaner_up 1
# HELP go_gc_cycles_automatic_gc_cycles_total Count of completed GC cycles generated by the Go runtime.
# TYPE go_gc_cycles_automatic_gc_cycles_total counter
go_gc_cycles_automatic_gc_cycles_total 0
# HELP go_gc_cycles_forced_gc_cycles_total Count of completed GC cycles forced by the application.
# TYPE go_gc_cycles_forced_gc_cycles_total counter
go_gc_cycles_forced_gc_cycles_total 0
# HELP go_gc_cycles_total_gc_cycles_total Count of all completed GC cycles.
# TYPE go_gc_cycles_total_gc_cycles_total counter
go_gc_cycles_total_gc_cycles_total 0
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0
go_gc_duration_seconds_sum 0
go_gc_duration_seconds_count 0
# HELP go_gc_heap_allocs_by_size_bytes_total Distribution of heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_allocs_by_size_bytes_total histogram
go_gc_heap_allocs_by_size_bytes_total_bucket{le="8.999999999999998"} 3072
go_gc_heap_allocs_by_size_bytes_total_bucket{le="24.999999999999996"} 13310
go_gc_heap_allocs_by_size_bytes_total_bucket{le="64.99999999999999"} 18168
go_gc_heap_allocs_by_size_bytes_total_bucket{le="144.99999999999997"} 25860
go_gc_heap_allocs_by_size_bytes_total_bucket{le="320.99999999999994"} 27495
go_gc_heap_allocs_by_size_bytes_total_bucket{le="704.9999999999999"} 28129
go_gc_heap_allocs_by_size_bytes_total_bucket{le="1536.9999999999998"} 28345
go_gc_heap_allocs_by_size_bytes_total_bucket{le="3200.9999999999995"} 28491
go_gc_heap_allocs_by_size_bytes_total_bucket{le="6528.999999999999"} 28552
go_gc_heap_allocs_by_size_bytes_total_bucket{le="13568.999999999998"} 28603
go_gc_heap_allocs_by_size_bytes_total_bucket{le="27264.999999999996"} 28613
go_gc_heap_allocs_by_size_bytes_total_bucket{le="+Inf"} 28615
go_gc_heap_allocs_by_size_bytes_total_sum 3.63168e+06
go_gc_heap_allocs_by_size_bytes_total_count 28615
# HELP go_gc_heap_allocs_bytes_total Cumulative sum of memory allocated to the heap by the application.
# TYPE go_gc_heap_allocs_bytes_total counter
go_gc_heap_allocs_bytes_total 3.63168e+06
# HELP go_gc_heap_allocs_objects_total Cumulative count of heap allocations triggered by the application. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_allocs_objects_total counter
go_gc_heap_allocs_objects_total 28615
# HELP go_gc_heap_frees_by_size_bytes_total Distribution of freed heap allocations by approximate size. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_frees_by_size_bytes_total histogram
go_gc_heap_frees_by_size_bytes_total_bucket{le="8.999999999999998"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="24.999999999999996"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="64.99999999999999"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="144.99999999999997"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="320.99999999999994"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="704.9999999999999"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="1536.9999999999998"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="3200.9999999999995"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="6528.999999999999"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="13568.999999999998"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="27264.999999999996"} 0
go_gc_heap_frees_by_size_bytes_total_bucket{le="+Inf"} 0
go_gc_heap_frees_by_size_bytes_total_sum 0
go_gc_heap_frees_by_size_bytes_total_count 0
# HELP go_gc_heap_frees_bytes_total Cumulative sum of heap memory freed by the garbage collector.
# TYPE go_gc_heap_frees_bytes_total counter
go_gc_heap_frees_bytes_total 0
# HELP go_gc_heap_frees_objects_total Cumulative count of heap allocations whose storage was freed by the garbage collector. Note that this does not include tiny objects as defined by /gc/heap/tiny/allocs:objects, only tiny blocks.
# TYPE go_gc_heap_frees_objects_total counter
go_gc_heap_frees_objects_total 0
# HELP go_gc_heap_goal_bytes Heap size target for the end of the GC cycle.
# TYPE go_gc_heap_goal_bytes gauge
go_gc_heap_goal_bytes 4.473924e+06
# HELP go_gc_heap_objects_objects Number of objects, live or unswept, occupying heap memory.
# TYPE go_gc_heap_objects_objects gauge
go_gc_heap_objects_objects 28615
# HELP go_gc_heap_tiny_allocs_objects_total Count of small allocations that are packed together into blocks. These allocations are counted separately from other allocations because each individual allocation is not tracked by the runtime, only their block. Each block is already accounted for in allocs-by-size and frees-by-size.
# TYPE go_gc_heap_tiny_allocs_objects_total counter
go_gc_heap_tiny_allocs_objects_total 674
# HELP go_gc_pauses_seconds_total Distribution individual GC-related stop-the-world pause latencies.
# TYPE go_gc_pauses_seconds_total histogram
go_gc_pauses_seconds_total_bucket{le="-5e-324"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-10"} 0
go_gc_pauses_seconds_total_bucket{le="9.999999999999999e-09"} 0
go_gc_pauses_seconds_total_bucket{le="1.2799999999999998e-07"} 0
go_gc_pauses_seconds_total_bucket{le="1.2799999999999998e-06"} 0
go_gc_pauses_seconds_total_bucket{le="1.6383999999999998e-05"} 0
go_gc_pauses_seconds_total_bucket{le="0.00016383999999999998"} 0
go_gc_pauses_seconds_total_bucket{le="0.0020971519999999997"} 0
go_gc_pauses_seconds_total_bucket{le="0.020971519999999997"} 0
go_gc_pauses_seconds_total_bucket{le="0.26843545599999996"} 0
go_gc_pauses_seconds_total_bucket{le="+Inf"} 0
go_gc_pauses_seconds_total_sum NaN
go_gc_pauses_seconds_total_count 0
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 9
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.17.7"} 1
# HELP go_memory_classes_heap_free_bytes Memory that is completely free and eligible to be returned to the underlying system, but has not been. This metric is the runtime's estimate of free address space that is backed by physical memory.
# TYPE go_memory_classes_heap_free_bytes gauge
go_memory_classes_heap_free_bytes 0
# HELP go_memory_classes_heap_objects_bytes Memory occupied by live objects and dead objects that have not yet been marked free by the garbage collector.
# TYPE go_memory_classes_heap_objects_bytes gauge
go_memory_classes_heap_objects_bytes 3.63168e+06
# HELP go_memory_classes_heap_released_bytes Memory that is completely free and has been returned to the underlying system. This metric is the runtime's estimate of free address space that is still mapped into the process, but is not backed by physical memory.
# TYPE go_memory_classes_heap_released_bytes gauge
go_memory_classes_heap_released_bytes 4.030464e+06
# HELP go_memory_classes_heap_stacks_bytes Memory allocated from the heap that is reserved for stack space, whether or not it is currently in-use.
# TYPE go_memory_classes_heap_stacks_bytes gauge
go_memory_classes_heap_stacks_bytes 688128
# HELP go_memory_classes_heap_unused_bytes Memory that is reserved for heap objects but is not currently used to hold heap objects.
# TYPE go_memory_classes_heap_unused_bytes gauge
go_memory_classes_heap_unused_bytes 38336
# HELP go_memory_classes_metadata_mcache_free_bytes Memory that is reserved for runtime mcache structures, but not in-use.
# TYPE go_memory_classes_metadata_mcache_free_bytes gauge
go_memory_classes_metadata_mcache_free_bytes 13568
# HELP go_memory_classes_metadata_mcache_inuse_bytes Memory that is occupied by runtime mcache structures that are currently being used.
# TYPE go_memory_classes_metadata_mcache_inuse_bytes gauge
go_memory_classes_metadata_mcache_inuse_bytes 19200
# HELP go_memory_classes_metadata_mspan_free_bytes Memory that is reserved for runtime mspan structures, but not in-use.
# TYPE go_memory_classes_metadata_mspan_free_bytes gauge
go_memory_classes_metadata_mspan_free_bytes 10184
# HELP go_memory_classes_metadata_mspan_inuse_bytes Memory that is occupied by runtime mspan structures that are currently being used.
# TYPE go_memory_classes_metadata_mspan_inuse_bytes gauge
go_memory_classes_metadata_mspan_inuse_bytes 55352
# HELP go_memory_classes_metadata_other_bytes Memory that is reserved for or used to hold runtime metadata.
# TYPE go_memory_classes_metadata_other_bytes gauge
go_memory_classes_metadata_other_bytes 3.647576e+06
# HELP go_memory_classes_os_stacks_bytes Stack memory allocated by the underlying operating system.
# TYPE go_memory_classes_os_stacks_bytes gauge
go_memory_classes_os_stacks_bytes 0
# HELP go_memory_classes_other_bytes Memory used by execution trace buffers, structures for debugging the runtime, finalizer and profiler specials, and more.
# TYPE go_memory_classes_other_bytes gauge
go_memory_classes_other_bytes 791300
# HELP go_memory_classes_profiling_buckets_bytes Memory that is used by the stack trace hash map used for profiling.
# TYPE go_memory_classes_profiling_buckets_bytes gauge
go_memory_classes_profiling_buckets_bytes 4276
# HELP go_memory_classes_total_bytes All memory mapped by the Go runtime into the current process as read-write. Note that this does not include memory mapped by code called via cgo or via the syscall package. Sum of all metrics in /memory/classes.
# TYPE go_memory_classes_total_bytes gauge
go_memory_classes_total_bytes 1.2930064e+07
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 3.63168e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 3.63168e+06
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 4276
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 674
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 0
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 3.647576e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 3.63168e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 4.030464e+06
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 3.670016e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 28615
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 4.030464e+06
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 7.70048e+06
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 0
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 29289
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 19200
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 32768
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 55352
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 65536
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 4.473924e+06
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 791300
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 688128
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 688128
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 1.2930064e+07
# HELP go_sched_goroutines_goroutines Count of live goroutines.
# TYPE go_sched_goroutines_goroutines gauge
go_sched_goroutines_goroutines 9
# HELP go_sched_latencies_seconds Distribution of the time goroutines have spent in the scheduler in a runnable state before actually running.
# TYPE go_sched_latencies_seconds histogram
go_sched_latencies_seconds_bucket{le="-5e-324"} 0
go_sched_latencies_seconds_bucket{le="9.999999999999999e-10"} 17
go_sched_latencies_seconds_bucket{le="9.999999999999999e-09"} 17
go_sched_latencies_seconds_bucket{le="1.2799999999999998e-07"} 17
go_sched_latencies_seconds_bucket{le="1.2799999999999998e-06"} 19
go_sched_latencies_seconds_bucket{le="1.6383999999999998e-05"} 22
go_sched_latencies_seconds_bucket{le="0.00016383999999999998"} 24
go_sched_latencies_seconds_bucket{le="0.0020971519999999997"} 24
go_sched_latencies_seconds_bucket{le="0.020971519999999997"} 24
go_sched_latencies_seconds_bucket{le="0.26843545599999996"} 24
go_sched_latencies_seconds_bucket{le="+Inf"} 24
go_sched_latencies_seconds_sum NaN
go_sched_latencies_seconds_count 24
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 7
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 0
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0

```

License
-------
Code is licensed under the [MIT](https://github.com/foxdalas/docker-cleaner/blob/master/LICENSE).
