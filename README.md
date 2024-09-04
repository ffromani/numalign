# numalign - tooling to check NUMA resource alignment/positioning

## DEPRECATION NOTICE

This package is declared obsolete and unmaintained.
[ctrreschk](https://github.com/ffromani/ctrreschk) is a more modern, featureful and actively development replacement.

## content

`irqcheck` tells information about IRQ/softirq cpus affinity.

`lsnt` reports information about NUMA locality of CPU and devices.

`numalign` tells you if a set of resources is aligned on the same NUMA node.

`sriovscan` finds all the SR-IOV (PFs and VFs) devices on the system and report infos about them.

`sriovctl` is a helper tool to override (kernel allowing) the NUMA placement of SRIOV devices, in case of buggy firmware.

`splitcpulist` parses a cpulist description and emits the list of all CPUs involved, to be used in shell scripts.

## license
(C) 2020 Red Hat Inc and licensed under the Apache License v2

## build
just run
```bash
make
```

See the READMEs under `cmd` for informations about each tool.

