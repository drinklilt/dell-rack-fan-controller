# Dell Rack Fan Controller
> A little piece of software to not deafen you

# Requirements
- [ipmitool](https://github.com/ipmitool/ipmitool)
- run as a user with access to /dev/ipmiX

# Useage
You can supply a path of the ipmi device with -device=/dev/ipmiX
```bash
dell-rack-fan-controller -device=/dev/ipmi0
```