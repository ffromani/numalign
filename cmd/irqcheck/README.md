## irqcheck

`irqcheck` tells about the IRQ and softirq cpus affinity

### Example output

IRQ:
```bash
IRQ   0 [                        ]: can run on [0 1 2 3]
IRQ   1 [                   i8042]: can run on [0 1 2 3]
IRQ   2 [                        ]: can run on [0 1 2 3]
IRQ   3 [                        ]: can run on [0 1 2 3]
IRQ   4 [                        ]: can run on [0 1 2 3]
IRQ   5 [                        ]: can run on [0 1 2 3]
IRQ   6 [                        ]: can run on [0 1 2 3]
IRQ   7 [                        ]: can run on [0 1 2 3]
IRQ   8 [                    rtc0]: can run on [0 1 2 3]
IRQ   9 [                    acpi]: can run on [0 1 2 3]
IRQ  10 [                        ]: can run on [0 1 2 3]
IRQ  11 [                        ]: can run on [0 1 2 3]
IRQ  12 [                   i8042]: can run on [0 1 2 3]
IRQ  13 [                        ]: can run on [0 1 2 3]
IRQ  14 [                        ]: can run on [0 1 2 3]
IRQ  15 [                        ]: can run on [0 1 2 3]
IRQ  16 [              i801_smbus]: can run on [0 1 2 3]
IRQ  19 [                        ]: can run on [0 1 2 3]
IRQ 120 [                        ]: can run on [0]
IRQ 121 [                        ]: can run on [0]
IRQ 125 [                xhci_hcd]: can run on [0 1 2 3]
IRQ 126 [               enp0s31f6]: can run on [0 1 2 3]
IRQ 127 [                 nvme0q0]: can run on [0 1 2 3]
IRQ 128 [                 nvme0q1]: can run on [0]
IRQ 129 [                 nvme0q2]: can run on [1]
IRQ 130 [                 nvme0q3]: can run on [2]
IRQ 131 [                 nvme0q4]: can run on [3]
IRQ 132 [                    i915]: can run on [0 1 2 3]
IRQ 133 [                  mei_me]: can run on [0 1 2 3]
IRQ 134 [     snd_hda_intel:card0]: can run on [0 1 2 3]
IRQ 135 [                 iwlwifi]: can run on [0 1 2 3]
IRQ 136 [              rmi4_smbus]: can run on [0 1 2 3]
IRQ 137 [            rmi4-00.fn34]: can run on [0 1 2 3]
IRQ 138 [            rmi4-00.fn01]: can run on [0 1 2 3]
IRQ 139 [            rmi4-00.fn03]: can run on [0 1 2 3]
IRQ 140 [            rmi4-00.fn11]: can run on [0 1 2 3]
IRQ 141 [            rmi4-00.fn11]: can run on [0 1 2 3]
IRQ 142 [            rmi4-00.fn30]: can run on [0 1 2 3]
```

softirq:
```bash
      HI = 0-3
   TIMER = 0-3
  NET_TX = 0-3
  NET_RX = 0-3
   BLOCK = 0-3
IRQ_POLL = 
 TASKLET = 0-3
   SCHED = 0-3
 HRTIMER = 3
     RCU = 0-3
```
