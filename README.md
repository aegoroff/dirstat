dirstat
=======

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/771c99225c8c4ad0a83115d2403f9fa0)](https://app.codacy.com/manual/egoroff/dirstat?utm_source=github.com&utm_medium=referral&utm_content=aegoroff/dirstat&utm_campaign=Badge_Grade_Dashboard) [![codecov](https://codecov.io/gh/aegoroff/dirstat/branch/master/graph/badge.svg)](https://codecov.io/gh/aegoroff/dirstat) [![Go Report Card](https://goreportcard.com/badge/github.com/aegoroff/dirstat)](https://goreportcard.com/report/github.com/aegoroff/dirstat)

Small tool that shows selected folder or drive (on Windows) usage statistic.
The tool shows file and folder statistic like the number of files or folders, top 10 files and folders by size,
file statistics by extension and by file size range.

## Install the pre-compiled binary

**homebrew** (only on macOS and Linux for now):

Add my tap (do it once):
```sh
brew tap aegoroff/tap
```
And then install dirstat:
```sh
brew install dirstat
```
Update dirstat if already installed:
```sh
brew upgrade dirstat
```

**scoop**:

```sh
scoop bucket add aegoroff https://github.com/aegoroff/scoop-bucket.git
scoop install dirstat
```

**manually**:

Download the pre-compiled binaries from the [releases](https://github.com/aegoroff/dirstat/releases) and
copy to the desired location.


Syntax:
------
```
 A small tool that shows selected folder or drive (on Windows) usage statistic

Usage:
  dirstat [flags]
  dirstat [command]

Available Commands:
  a           Show all information about folder/volume
  b           Show the first digit distribution of files size (benford law validation)
  e           Show file extensions statistic
  fi          Show information only about files
  fo          Show information only about folders
  help        Help about any command
  version     Print the version number of dirstat

Flags:
  -h, --help            help for dirstat
  -m, --memory          Show memory statistic after run
  -p, --output string   Write results into file. Specify path to output file using this option
  -o, --removeroot      Remove root part from full path i.e. output relative paths
  -t, --top int         The number of lines in top statistics. (default 10)

Use "dirstat [command] --help" for more information about a command.
```
Examples:
---------
Show statistic about all D drive (on Windows)
```
dirstat a d:
```
Show all statistic about specific folder d:\folder
```
dirstat a d:\folder
```
Show all statistic and additional verbose statistic about files in ranges Between 10 GiB and 100 GiB and Between 100 GiB and 1.0 TiB 
```
dirstat a d:\ -r 7 -r 8
```
or
```
dirstat a d:\ -r 7,8
```
The second form is equivalent

Output example:
---------------
```
Root: c:\

Total files stat:

┌────┬─────────────────────────────┬─────────┬────────┬────────┬────────┐
│  # │ FILE SIZE                   │ AMOUNT  │ %      │ SIZE   │ %      │
├────┼─────────────────────────────┼─────────┼────────┼────────┼────────┤
│  1 │ Between 0 B and 100 KiB     │ 1432007 │ 88.09% │ 15 GiB │ 4.78%  │
│  2 │ Between 100 KiB and 1.0 MiB │ 157685  │ 9.70%  │ 47 GiB │ 14.81% │
│  3 │ Between 1.0 MiB and 10 MiB  │ 32480   │ 2.00%  │ 83 GiB │ 26.00% │
│  4 │ Between 10 MiB and 100 MiB  │ 3731    │ 0.23%  │ 95 GiB │ 29.73% │
│  5 │ Between 100 MiB and 1.0 GiB │ 175     │ 0.01%  │ 37 GiB │ 11.57% │
│  6 │ Between 1.0 GiB and 10 GiB  │ 4       │ 0.00%  │ 18 GiB │ 5.60%  │
│  7 │ Between 10 GiB and 100 GiB  │ 1       │ 0.00%  │ 24 GiB │ 7.70%  │
│  8 │ Between 100 GiB and 1.0 TiB │ 0       │ 0.00%  │ 0 B    │ 0.00%  │
│  9 │ Between 1.0 TiB and 10 TiB  │ 0       │ 0.00%  │ 0 B    │ 0.00%  │
│ 10 │ Between 10 TiB and 1.0 PiB  │ 0       │ 0.00%  │ 0 B    │ 0.00%  │
└────┴─────────────────────────────┴─────────┴────────┴────────┴────────┘

TOP 10 file extensions by size:

┌────┬───────────┬────────┬────────┬─────────┬────────┐
│  # │ EXTENSION │ COUNT  │ %      │ SIZE    │ %      │
├────┼───────────┼────────┼────────┼─────────┼────────┤
│  1 │ .dll      │ 193703 │ 11.92% │ 77 GiB  │ 24.20% │
│  2 │ .vhdx     │ 5      │ 0.00%  │ 37 GiB  │ 11.65% │
│  3 │ .exe      │ 18009  │ 1.11%  │ 19 GiB  │ 6.03%  │
│  4 │ .lib      │ 15445  │ 0.95%  │ 18 GiB  │ 5.61%  │
│  5 │           │ 145504 │ 8.95%  │ 13 GiB  │ 4.20%  │
│  6 │ .jar      │ 9786   │ 0.60%  │ 12 GiB  │ 3.82%  │
│  7 │ .cab      │ 2211   │ 0.14%  │ 11 GiB  │ 3.49%  │
│  8 │ .msi      │ 1998   │ 0.12%  │ 11 GiB  │ 3.44%  │
│  9 │ .xml      │ 67058  │ 4.13%  │ 10 GiB  │ 3.24%  │
│ 10 │ .pdb      │ 8963   │ 0.55%  │ 6.6 GiB │ 2.06%  │
└────┴───────────┴────────┴────────┴─────────┴────────┘

TOP 10 file extensions by count:

┌────┬───────────┬────────┬────────┬─────────┬────────┐
│  # │ EXTENSION │ COUNT  │ %      │ SIZE    │ %      │
├────┼───────────┼────────┼────────┼─────────┼────────┤
│  1 │ .dll      │ 193703 │ 11.92% │ 77 GiB  │ 24.20% │
│  2 │           │ 145504 │ 8.95%  │ 13 GiB  │ 4.20%  │
│  3 │ .py       │ 74963  │ 4.61%  │ 766 MiB │ 0.23%  │
│  4 │ .js       │ 69719  │ 4.29%  │ 2.3 GiB │ 0.73%  │
│  5 │ .xml      │ 67058  │ 4.13%  │ 10 GiB  │ 3.24%  │
│  6 │ .png      │ 61267  │ 3.77%  │ 868 MiB │ 0.27%  │
│  7 │ .h        │ 60334  │ 3.71%  │ 2.5 GiB │ 0.79%  │
│  8 │ .pyc      │ 50239  │ 3.09%  │ 525 MiB │ 0.16%  │
│  9 │ .manifest │ 47668  │ 2.93%  │ 380 MiB │ 0.12%  │
│ 10 │ .html     │ 46017  │ 2.83%  │ 616 MiB │ 0.19%  │
└────┴───────────┴────────┴────────┴─────────┴────────┘

TOP 10 files by size:

┌────┬──────────────────────────────────────────────────────────────────────────────────────────────────────┬─────────┐
│  # │ FILE                                                                                                 │ SIZE    │
├────┼──────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────┤
│  1 │ c:\Users\egr\AppData\Local\Docker\wsl\data\ext4.vhdx                                                 │ 24 GiB  │
│  2 │ c:\ProgramData\DockerDesktop\vm-data\DockerDesktop.vhdx                                              │ 9.5 GiB │
│  3 │ c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.Ubuntu20.04onWindows_79rhkp1fndgsc\LocalSt │ 2.9 GiB │
│    │ ate\ext4.vhdx                                                                                        │         │
│  4 │ c:\Program Files\MongoDB\Server\4.0\log\mongod.log                                                   │ 2.9 GiB │
│  5 │ c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0\flash.vhd                       │ 2.6 GiB │
│  6 │ c:\Program Files\MongoDB\Server\4.0\data\collection-7--2807253782407563390.wt                        │ 962 MiB │
│  7 │ c:\Program Files\MongoDB\Server\4.0\data\collection-2-2535253104322579421.wt                         │ 946 MiB │
│  8 │ c:\Windows\Installer\13b2694.msi                                                                     │ 887 MiB │
│  9 │ c:\Windows\Installer\9da85.msp                                                                       │ 779 MiB │
│ 10 │ c:\Program Files\NVIDIA GPU Computing Toolkit\CUDA\v11.1\bin\cudnn_cnn_infer64_8.dll                 │ 673 MiB │
└────┴──────────────────────────────────────────────────────────────────────────────────────────────────────┴─────────┘

TOP 10 folders by size:

┌────┬──────────────────────────────────────────────────────────────────────────────────────────────────────┬───────┬───────┬─────────┬───────┐
│  # │ FOLDER                                                                                               │ FILES │ %     │ SIZE    │ %     │
├────┼──────────────────────────────────────────────────────────────────────────────────────────────────────┼───────┼───────┼─────────┼───────┤
│  1 │ c:\Users\egr\AppData\Local\Docker\wsl\data                                                           │ 1     │ 0.00% │ 24 GiB  │ 7.70% │
│  2 │ c:\Windows\Installer                                                                                 │ 1976  │ 0.12% │ 11 GiB  │ 3.54% │
│  3 │ c:\ProgramData\DockerDesktop\vm-data                                                                 │ 1     │ 0.00% │ 9.5 GiB │ 2.97% │
│  4 │ c:\Program Files\MongoDB\Server\4.0\data                                                             │ 41    │ 0.00% │ 3.5 GiB │ 1.08% │
│  5 │ c:\Program Files\NVIDIA GPU Computing Toolkit\CUDA\v11.1\bin                                         │ 47    │ 0.00% │ 3.4 GiB │ 1.06% │
│  6 │ c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.Ubuntu20.04onWindows_79rhkp1fndgsc\LocalSt │ 1     │ 0.00% │ 2.9 GiB │ 0.92% │
│    │ ate                                                                                                  │       │       │         │       │
│  7 │ c:\Program Files\MongoDB\Server\4.0\log                                                              │ 1     │ 0.00% │ 2.9 GiB │ 0.91% │
│  8 │ c:\Windows\System32                                                                                  │ 5053  │ 0.31% │ 2.6 GiB │ 0.81% │
│  9 │ c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0                                 │ 1     │ 0.00% │ 2.6 GiB │ 0.80% │
│ 10 │ c:\ProgramData\Package Cache\{EA923538-8370-4294-A5CC-F6130FAAD89D}v10.1.10586.11\Redistributable\1. │ 51    │ 0.00% │ 1.5 GiB │ 0.47% │
│    │ 0.0.0                                                                                                │       │       │         │       │
└────┴──────────────────────────────────────────────────────────────────────────────────────────────────────┴───────┴───────┴─────────┴───────┘

TOP 10 folders by count:

┌────┬─────────────────────────────────────────────────────────────────────────────────────────────────┬───────┬───────┬─────────┬───────┐
│  # │ FOLDER                                                                                          │ FILES │ %     │ SIZE    │ %     │
├────┼─────────────────────────────────────────────────────────────────────────────────────────────────┼───────┼───────┼─────────┼───────┤
│  1 │ c:\Windows\servicing\LCU\Package_for_RollupFix~31bf3856ad364e35~amd64~~19041.685.1.6            │ 32058 │ 1.97% │ 492 MiB │ 0.15% │
│  2 │ c:\Windows\servicing\LCU\Package_for_RollupFix~31bf3856ad364e35~amd64~~19041.746.1.6            │ 31644 │ 1.95% │ 472 MiB │ 0.14% │
│  3 │ c:\Windows\WinSxS\Manifests                                                                     │ 25373 │ 1.56% │ 23 MiB  │ 0.01% │
│  4 │ c:\Users\egr\AppData\Local\Google\Chrome\User Data\Default\Local Storage                        │ 6862  │ 0.42% │ 84 MiB  │ 0.03% │
│  5 │ c:\msys64\usr\share\man\man3                                                                    │ 6739  │ 0.41% │ 24 MiB  │ 0.01% │
│  6 │ c:\Users\egr\AppData\Roaming\Microsoft\Crypto\RSA\S-1-5-21-2126183571-1528345743-823968152-1001 │ 5144  │ 0.32% │ 15 MiB  │ 0.00% │
│  7 │ c:\Windows\System32                                                                             │ 5053  │ 0.31% │ 2.6 GiB │ 0.81% │
│  8 │ c:\Windows\servicing\Packages                                                                   │ 4490  │ 0.28% │ 62 MiB  │ 0.02% │
│  9 │ c:\msys64\usr\share\doc\openssl\html\man3                                                       │ 3792  │ 0.23% │ 36 MiB  │ 0.01% │
│ 10 │ c:\Users\egr\AppData\Local\Microsoft\Edge\User Data\Default\Code Cache\js                       │ 3771  │ 0.23% │ 286 MiB │ 0.09% │
└────┴─────────────────────────────────────────────────────────────────────────────────────────────────┴───────┴───────┴─────────┴───────┘

Total files:            1625602 (319 GiB)
Total folders:          365894
Total file extensions:  6401

Read taken:	5.9741123s

Alloc = 3.3 MiB	TotalAlloc = 2.7 GiB	Sys = 82 MiB	NumGC = 453	NumGoRoutines = 1
```