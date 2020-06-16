# numalign - tooling to check NUMA resource alignment/positioning

`lsnt` reports information about NUMA locality of CPU and devices.

`numalign` tells you if a set of resources is aligned on the same NUMA node.

`sriovscan` finds all the SR-IOV (PFs and VFs) devices on the system and report infos about them.

## license
(C) 2020 Red Hat Inc and licensed under the Apache License v2

## build
just run
```bash
make
```

## lsnt

TBD

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


## numalign

### Container image
```bash
podman run -e NUMALIGN_DEBUG=1 quay.io/fromani/numalign:devel
```

### Example output
From a developer laptop:
```bash
$ ./numalign 
$
$ # no output, let's see why
$ NUMALIGN_DEBUG=1 ./numalign
2020/01/24 14:00:41 CPU: allowed: [0 1 2 3]
2020/01/24 14:00:41 CPU: NUMA node by id: map[0:0 1:0 2:0 3:0]
2020/01/24 14:00:41 No PCI devices detected
$
$ # OK, let's add some PCI devices to be checked. You need to use valid PCI ids
$ # (e.g referring to actual devices on the system on which you are running `numalign`)
$ # use `lspci` to find some.
$ #
$ # past the PCIDEVICE_ prefix, the rest of the variable name is not really important
$ export PCIDEVICE_FOOBAR="0000:3c:00.0"
$
$ NUMALIGN_DEBUG=1 ./numalign 
2020/01/24 14:02:50 CPU: allowed: [0 1 2 3]
2020/01/24 14:02:50 CPU: NUMA node by id: map[0:0 1:0 2:0 3:0]
2020/01/24 14:02:50 PCI: devices: 0000:3c:00.0
STATUS ALIGNED=true
NUMA NODE=0
^C
$
$ # what about a quick(er) check?
$ NUMALIGN_DEBUG=1 NUMALIGN_SLEEP_HOURS=0 ./numalign 
2020/01/27 14:32:54 CPU: allowed: [0 1 2 3]
2020/01/27 14:32:54 CPU: NUMA node by id: map[0:0 1:0 2:0 3:0]
2020/01/27 14:32:54 PCI: devices: 0000:3c:00.0
STATUS ALIGNED=true
NUMA NODE=0
$
$ # you can also check for other processes in the same container
$ NUMALIGN_DEBUG=1 NUMALIGN_SLEEP_HOURS=0 NUMALIGN_DEBUG=1 ./numalign 1700
2020/05/21 10:48:25 CPU: allowed for "1700": [0 1 2 3]
2020/05/21 10:48:25 CPU: NUMA node by id: map[0:0 1:0 2:0 3:0]
2020/05/21 10:48:25 PCI: devices: 0000:3c:00.0
STATUS ALIGNED=true
NUMA NODE=0
$
$ NUMALIGN_DEBUG=1 NUMALIGN_SLEEP_HOURS=0 ./numalign 1700 self
2020/05/21 10:48:36 CPU: allowed for "1700": [0 1 2 3]
2020/05/21 10:48:36 CPU: allowed for "self": [0 1 2 3]
2020/05/21 10:48:36 CPU: NUMA node by id: map[0:0 1:0 2:0 3:0]
2020/05/21 10:48:36 PCI: devices: 0000:3c:00.0
STATUS ALIGNED=true
NUMA NODE=0
$
```

## sriovscan

TBD
