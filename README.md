# numalign - a simple tool to check resource alignment.

`numalign` tells you if a set of resources is aligned on the same NUMA node. That's it.

## license
numalign (C) 2020 Red Hat Inc and licensed under the Apache License v2

## build
just run
```bash
make
```

## Container image
```bash
podman run -e NUMALIGN_DEBUG=1 quay.io/fromani/numalign:devel
```

## Example output
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

```
