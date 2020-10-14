## lsnt

`lsnt` reports information regarding NUMA Topology, like CPU placement, PCI devices placement.

### Example output
On a multi-NUMA system:
```bash
$ lsnt
Usage:
  lsnt [flags]
  lsnt [command]

Available Commands:
  cpu         show cpu details like lscpu(1)
  daemonwait  wait forever, or until a UNIX signal (SIGINT, SIGTERM) arrives
  help        Help about any command
  numa        show NUMA device tree
  pcidevs     show PCI devices in the system

Flags:
  -h, --help           help for lsnt
  -S, --sysfs string   sysfs root (default "/sys")
      --verbose int    verbosiness level (default 1)

Use "lsnt [command] --help" for more information about a command.
$
$
$ # output intentionally similar to `lscpu`
$ lsnt cpu
CPU(s):              24
Present CPU(s) list: 0-23
On-line CPU(s) list: 0-23
Thread(s) per core:  2
Core(s) per socket:  12
Socket(s):           2
NUMA node(s):        2
NUMA node0 CPU(s):   0,2,4,6,8,10,12,14,16,18,20,22
NUMA node1 CPU(s):   1,3,5,7,9,11,13,15,17,19,21,23
$
$ # let's see it from another perspective
$ lsnt numa
.
└── numa00
│   ├── 0,2,4,6,8,10,12,14,16,18,20,22
└── numa01
    └── 1,3,5,7,9,11,13,15,17,19,21,23

$
$ # now the PCI devices:
$ lsnt pcidevs -N -T
.
└── UNKNOWN
    └── 0000:01:00.0 14e4:1639 (0200)
    └── 0000:01:00.1 14e4:1639 (0200)
    └── 0000:02:00.0 14e4:1639 (0200)
    └── 0000:02:00.1 14e4:1639 (0200)
    └── 0000:05:00.0 8086:1521 physfn numvfs=4
    └── 0000:05:00.1 8086:1521 physfn numvfs=4
    └── 0000:05:10.0 8086:1520 vfn parent=0000:05:00.0
    └── 0000:05:10.1 8086:1520 vfn parent=0000:05:00.1
    └── 0000:05:10.4 8086:1520 vfn parent=0000:05:00.0
    └── 0000:05:10.5 8086:1520 vfn parent=0000:05:00.1
    └── 0000:05:11.0 8086:1520 vfn parent=0000:05:00.0
    └── 0000:05:11.1 8086:1520 vfn parent=0000:05:00.1
    └── 0000:05:11.4 8086:1520 vfn parent=0000:05:00.0
    └── 0000:05:11.5 8086:1520 vfn parent=0000:05:00.1

$
$ # node reported UNKNOWN because of
$ # https://access.redhat.com/solutions/435313
$
$ # let's see the same information from another perspective
$ lsnt pcidevs -N -T -P
.
└── UNKNOWN
    └── 0000:01:00.0 14e4:1639 (0200)
    └── 0000:01:00.1 14e4:1639 (0200)
    └── 0000:02:00.0 14e4:1639 (0200)
    └── 0000:02:00.1 14e4:1639 (0200)
    └── 0000:05:00.0 8086:1521 physfn numvfs=4
    │   ├── 0000:05:10.0 8086:1520 vfn
    │   ├── 0000:05:10.4 8086:1520 vfn
    │   ├── 0000:05:11.0 8086:1520 vfn
    │   ├── 0000:05:11.4 8086:1520 vfn
    └── 0000:05:00.1 8086:1521 physfn numvfs=4
        └── 0000:05:10.1 8086:1520 vfn
        └── 0000:05:10.5 8086:1520 vfn
        └── 0000:05:11.1 8086:1520 vfn
        └── 0000:05:11.5 8086:1520 vfn


```

