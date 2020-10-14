## numalign

Simple NUMA Alignment check, to be used inside containers.

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

