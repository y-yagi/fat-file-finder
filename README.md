# fat-file-finder

`fat-file-finder` is a tool for finding large files.

## Installation

Use `go get` to install this package:

```bash
$ go get github.com/y-yagi/fat-file-finder
```

## Usage

```shell
fat-file-finder --help
Usage of ./fat-file-finder:
  -l string
    	Search location. (default ".")
  -s string
    	Threshold size to display. (default "100M")
  -v	show version
```

Example.

```
$ fat-file-finder -s 300M
f .AndroidStudioPreview2.0/system/index/android.value.resources.index/android.value.resources.index.values (427.6M)
f .android/avd/Nexus_5_API_22_x86.avd/userdata-qemu.img (550M)
f .android/avd/Nexus_5_API_22_x86.avd/userdata.img (550M)
f .android/avd/Nexus_S_Edited_API_19.avd/userdata-qemu.img (550M)
f .android/avd/Nexus_S_Edited_API_19.avd/userdata.img (550M
...
```

## Contribution

1. Fork ([https://github.com/y-yagi/fat-file-finder/fork](https://github.com/y-yagi/fat-file-finder/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
