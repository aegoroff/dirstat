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

┌────┬─────────────────────────────┬───────────┬────────┬────────┬────────┐
│  # │ FILE SIZE                   │ AMOUNT    │ %      │ SIZE   │ %      │
├────┼─────────────────────────────┼───────────┼────────┼────────┼────────┤
│  1 │ Between 0 B and 100 KiB     │ 1,422,801 │ 88.09% │ 15 GiB │ 4.76%  │
│  2 │ Between 100 KiB and 1.0 MiB │ 156,536   │ 9.69%  │ 47 GiB │ 14.81% │
│  3 │ Between 1.0 MiB and 10 MiB  │ 32,018    │ 1.98%  │ 82 GiB │ 25.94% │
│  4 │ Between 10 MiB and 100 MiB  │ 3,721     │ 0.23%  │ 94 GiB │ 29.63% │
│  5 │ Between 100 MiB and 1.0 GiB │ 172       │ 0.01%  │ 36 GiB │ 11.48% │
│  6 │ Between 1.0 GiB and 10 GiB  │ 4         │ 0.00%  │ 18 GiB │ 5.64%  │
│  7 │ Between 10 GiB and 100 GiB  │ 1         │ 0.00%  │ 24 GiB │ 7.75%  │
│  8 │ Between 100 GiB and 1.0 TiB │ 0         │ 0.00%  │ 0 B    │ 0.00%  │
│  9 │ Between 1.0 TiB and 10 TiB  │ 0         │ 0.00%  │ 0 B    │ 0.00%  │
│ 10 │ Between 10 TiB and 1.0 PiB  │ 0         │ 0.00%  │ 0 B    │ 0.00%  │
└────┴─────────────────────────────┴───────────┴────────┴────────┴────────┘

TOP 10 file extensions by size:

┌────┬───────────┬─────────┬────────┬─────────┬────────┐
│  # │ EXTENSION │ COUNT   │ %      │ SIZE    │ %      │
├────┼───────────┼─────────┼────────┼─────────┼────────┤
│  1 │ .dll      │ 187,010 │ 11.58% │ 76 GiB  │ 24.10% │
│  2 │ .vhdx     │ 5       │ 0.00%  │ 37 GiB  │ 11.72% │
│  3 │ .exe      │ 17,867  │ 1.11%  │ 19 GiB  │ 6.06%  │
│  4 │ .lib      │ 15,465  │ 0.96%  │ 18 GiB  │ 5.65%  │
│  5 │           │ 143,385 │ 8.88%  │ 13 GiB  │ 4.16%  │
│  6 │ .jar      │ 9,786   │ 0.61%  │ 12 GiB  │ 3.84%  │
│  7 │ .msi      │ 2,010   │ 0.12%  │ 12 GiB  │ 3.64%  │
│  8 │ .xml      │ 66,852  │ 4.14%  │ 10 GiB  │ 3.25%  │
│  9 │ .cab      │ 2,189   │ 0.14%  │ 10 GiB  │ 3.22%  │
│ 10 │ .pdb      │ 8,923   │ 0.55%  │ 6.6 GiB │ 2.08%  │
└────┴───────────┴─────────┴────────┴─────────┴────────┘

TOP 10 file extensions by count:

┌────┬───────────┬─────────┬────────┬─────────┬────────┐
│  # │ EXTENSION │ COUNT   │ %      │ SIZE    │ %      │
├────┼───────────┼─────────┼────────┼─────────┼────────┤
│  1 │ .dll      │ 187,010 │ 11.58% │ 76 GiB  │ 24.10% │
│  2 │           │ 143,385 │ 8.88%  │ 13 GiB  │ 4.16%  │
│  3 │ .py       │ 74,963  │ 4.64%  │ 766 MiB │ 0.24%  │
│  4 │ .js       │ 69,742  │ 4.32%  │ 2.3 GiB │ 0.74%  │
│  5 │ .xml      │ 66,852  │ 4.14%  │ 10 GiB  │ 3.25%  │
│  6 │ .png      │ 61,279  │ 3.79%  │ 868 MiB │ 0.27%  │
│  7 │ .h        │ 60,354  │ 3.74%  │ 2.5 GiB │ 0.80%  │
│  8 │ .pyc      │ 50,239  │ 3.11%  │ 525 MiB │ 0.16%  │
│  9 │ .manifest │ 47,652  │ 2.95%  │ 380 MiB │ 0.12%  │
│ 10 │ .html     │ 46,017  │ 2.85%  │ 616 MiB │ 0.19%  │
└────┴───────────┴─────────┴────────┴─────────┴────────┘

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
│  1 │ c:\Users\egr\AppData\Local\Docker\wsl\data                                                           │ 1     │ 0.00% │ 24 GiB  │ 7.75% │
│  2 │ c:\Windows\Installer                                                                                 │ 1,988 │ 0.12% │ 12 GiB  │ 3.65% │
│  3 │ c:\ProgramData\DockerDesktop\vm-data                                                                 │ 1     │ 0.00% │ 9.5 GiB │ 2.99% │
│  4 │ c:\Program Files\MongoDB\Server\4.0\data                                                             │ 41    │ 0.00% │ 3.5 GiB │ 1.09% │
│  5 │ c:\Program Files\NVIDIA GPU Computing Toolkit\CUDA\v11.1\bin                                         │ 47    │ 0.00% │ 3.4 GiB │ 1.07% │
│  6 │ c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.Ubuntu20.04onWindows_79rhkp1fndgsc\LocalSt │ 1     │ 0.00% │ 2.9 GiB │ 0.93% │
│    │ ate                                                                                                  │       │       │         │       │
│  7 │ c:\Program Files\MongoDB\Server\4.0\log                                                              │ 1     │ 0.00% │ 2.9 GiB │ 0.92% │
│  8 │ c:\Windows\System32                                                                                  │ 5,053 │ 0.31% │ 2.6 GiB │ 0.81% │
│  9 │ c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0                                 │ 1     │ 0.00% │ 2.6 GiB │ 0.81% │
│ 10 │ c:\ProgramData\Package Cache\{EA923538-8370-4294-A5CC-F6130FAAD89D}v10.1.10586.11\Redistributable\1. │ 51    │ 0.00% │ 1.5 GiB │ 0.48% │
│    │ 0.0.0                                                                                                │       │       │         │       │
└────┴──────────────────────────────────────────────────────────────────────────────────────────────────────┴───────┴───────┴─────────┴───────┘

TOP 10 folders by count:

┌────┬─────────────────────────────────────────────────────────────────────────────────────────────────┬────────┬───────┬─────────┬───────┐
│  # │ FOLDER                                                                                          │ FILES  │ %     │ SIZE    │ %     │
├────┼─────────────────────────────────────────────────────────────────────────────────────────────────┼────────┼───────┼─────────┼───────┤
│  1 │ c:\Windows\servicing\LCU\Package_for_RollupFix~31bf3856ad364e35~amd64~~19041.685.1.6            │ 32,058 │ 1.98% │ 492 MiB │ 0.15% │
│  2 │ c:\Windows\servicing\LCU\Package_for_RollupFix~31bf3856ad364e35~amd64~~19041.746.1.6            │ 31,644 │ 1.96% │ 472 MiB │ 0.15% │
│  3 │ c:\Windows\WinSxS\Manifests                                                                     │ 25,373 │ 1.57% │ 23 MiB  │ 0.01% │
│  4 │ c:\Users\egr\AppData\Local\Google\Chrome\User Data\Default\Local Storage                        │ 6,862  │ 0.42% │ 84 MiB  │ 0.03% │
│  5 │ c:\msys64\usr\share\man\man3                                                                    │ 6,739  │ 0.42% │ 24 MiB  │ 0.01% │
│  6 │ c:\Users\egr\AppData\Roaming\Microsoft\Crypto\RSA\S-1-5-21-2126183571-1528345743-823968152-1001 │ 5,144  │ 0.32% │ 15 MiB  │ 0.00% │
│  7 │ c:\Windows\System32                                                                             │ 5,053  │ 0.31% │ 2.6 GiB │ 0.81% │
│  8 │ c:\Windows\servicing\Packages                                                                   │ 4,490  │ 0.28% │ 62 MiB  │ 0.02% │
│  9 │ c:\msys64\usr\share\doc\openssl\html\man3                                                       │ 3,792  │ 0.23% │ 36 MiB  │ 0.01% │
│ 10 │ c:\Program Files\Adobe\Adobe Premiere Rush\PNG                                                  │ 3,587  │ 0.22% │ 15 MiB  │ 0.00% │
└────┴─────────────────────────────────────────────────────────────────────────────────────────────────┴────────┴───────┴─────────┴───────┘

Total files:            1,615,253 (317 GiB)
Total folders:          363,476
Total file extensions:  6,402

Read taken:	6.5079993s
```