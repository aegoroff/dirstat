# dirstat

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/771c99225c8c4ad0a83115d2403f9fa0)](https://app.codacy.com/manual/egoroff/dirstat?utm_source=github.com&utm_medium=referral&utm_content=aegoroff/dirstat&utm_campaign=Badge_Grade_Dashboard) [![codecov](https://codecov.io/gh/aegoroff/dirstat/branch/master/graph/badge.svg)](https://codecov.io/gh/aegoroff/dirstat)
[![Go Report Card](https://goreportcard.com/badge/github.com/aegoroff/dirstat)](https://goreportcard.com/report/github.com/aegoroff/dirstat)
[![CI](https://github.com/aegoroff/dirstat/actions/workflows/ci.yml/badge.svg)](https://github.com/aegoroff/dirstat/actions/workflows/ci.yml)
[![](https://tokei.rs/b1/github/aegoroff/dirstat?category=code)](https://github.com/XAMPPRocky/tokei)

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

**AUR (Arch Linux User Repository)**:

install binary package:

```sh
 yay -S dirstat-bin
```

or if yay reports that package not found force updating repo info

```sh
yay -Syyu dirstat-bin
```

**manually**:

Download the pre-compiled binaries from the [releases](https://github.com/aegoroff/dirstat/releases) and
copy to the desired location.

**install deb package on Arch Linux**:

1. Install [debtap](https://github.com/helixarch/debtap) from AUR using yay:

```sh
 yay -S debtap
```

2. Create equivalent package using debtap:

```sh
 sudo debtap -u
 debtap dirstat_x.x.x_amd64.deb
```

3. Install using pacman:

```sh
sudo pacman -U dirstat-x.x.x-1-x86_64.pkg.tar.zst
```

## Syntax:

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

## Examples:

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

## Output example:

Windows:

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

Mac OSX

```
Root: /System/Volumes/Data

Total files stat:

┌────┬─────────────────────────────┬─────────┬────────┬─────────┬────────┐
│  # │ FILE SIZE                   │ AMOUNT  │ %      │ SIZE    │ %      │
├────┼─────────────────────────────┼─────────┼────────┼─────────┼────────┤
│  1 │ Between 0 B and 100 KiB     │ 889,157 │ 93.72% │ 6.2 GiB │ 7.73%  │
│  2 │ Between 100 KiB and 1.0 MiB │ 49,683  │ 5.24%  │ 14 GiB  │ 17.26% │
│  3 │ Between 1.0 MiB and 10 MiB  │ 8,812   │ 0.93%  │ 23 GiB  │ 28.24% │
│  4 │ Between 10 MiB and 100 MiB  │ 1,017   │ 0.11%  │ 24 GiB  │ 29.36% │
│  5 │ Between 100 MiB and 1.0 GiB │ 79      │ 0.01%  │ 14 GiB  │ 17.41% │
│  6 │ Between 1.0 GiB and 10 GiB  │ 0       │ 0.00%  │ 0 B     │ 0.00%  │
│  7 │ Between 10 GiB and 100 GiB  │ 0       │ 0.00%  │ 0 B     │ 0.00%  │
│  8 │ Between 100 GiB and 1.0 TiB │ 0       │ 0.00%  │ 0 B     │ 0.00%  │
│  9 │ Between 1.0 TiB and 10 TiB  │ 0       │ 0.00%  │ 0 B     │ 0.00%  │
│ 10 │ Between 10 TiB and 1.0 PiB  │ 0       │ 0.00%  │ 0 B     │ 0.00%  │
└────┴─────────────────────────────┴─────────┴────────┴─────────┴────────┘

TOP 10 file extensions by size:

┌────┬───────────┬─────────┬────────┬─────────┬────────┐
│  # │ EXTENSION │ COUNT   │ %      │ SIZE    │ %      │
├────┼───────────┼─────────┼────────┼─────────┼────────┤
│  1 │           │ 126,747 │ 13.36% │ 23 GiB  │ 28.52% │
│  2 │ .jar      │ 4,287   │ 0.45%  │ 4.9 GiB │ 6.05%  │
│  3 │ .pdf      │ 1,695   │ 0.18%  │ 3.6 GiB │ 4.41%  │
│  4 │ .dylib    │ 1,795   │ 0.19%  │ 3.2 GiB │ 3.98%  │
│  5 │ .jpg      │ 5,555   │ 0.59%  │ 2.5 GiB │ 3.15%  │
│  6 │ .mp4      │ 110     │ 0.01%  │ 2.4 GiB │ 2.92%  │
│  7 │ .dmg      │ 9       │ 0.00%  │ 2.0 GiB │ 2.42%  │
│  8 │ .caf      │ 2,212   │ 0.23%  │ 1.6 GiB │ 1.96%  │
│  9 │ .ttc      │ 136     │ 0.01%  │ 1.5 GiB │ 1.90%  │
│ 10 │ .dll      │ 8,409   │ 0.89%  │ 1.5 GiB │ 1.89%  │
└────┴───────────┴─────────┴────────┴─────────┴────────┘

TOP 10 file extensions by count:

┌────┬───────────┬─────────┬────────┬─────────┬────────┐
│  # │ EXTENSION │ COUNT   │ %      │ SIZE    │ %      │
├────┼───────────┼─────────┼────────┼─────────┼────────┤
│  1 │ .strings  │ 220,260 │ 23.22% │ 1.2 GiB │ 1.43%  │
│  2 │           │ 126,747 │ 13.36% │ 23 GiB  │ 28.52% │
│  3 │ .js       │ 68,841  │ 7.26%  │ 1.2 GiB │ 1.44%  │
│  4 │ .go       │ 51,202  │ 5.40%  │ 852 MiB │ 1.03%  │
│  5 │ .png      │ 49,237  │ 5.19%  │ 894 MiB │ 1.08%  │
│  6 │ .html     │ 43,987  │ 4.64%  │ 229 MiB │ 0.28%  │
│  7 │ .h        │ 23,741  │ 2.50%  │ 274 MiB │ 0.33%  │
│  8 │ .3pm      │ 18,018  │ 1.90%  │ 234 MiB │ 0.28%  │
│  9 │ .plist    │ 17,127  │ 1.81%  │ 253 MiB │ 0.31%  │
│ 10 │ .json     │ 16,961  │ 1.79%  │ 455 MiB │ 0.55%  │
└────┴───────────┴─────────┴────────┴─────────┴────────┘

TOP 10 files by size:

┌────┬──────────────────────────────────────────────────────────────────────────────────────────────────────┬─────────┐
│  # │ FILE                                                                                                 │ SIZE    │
├────┼──────────────────────────────────────────────────────────────────────────────────────────────────────┼─────────┤
│  1 │ /private/var/vm/sleepimage                                                                           │ 1.0 GiB │
│  2 │ /System/Library/AssetsV2/com_apple_MobileAsset_MacSoftwareUpdate/422e7e556d3f01a4bbf259b0f10f7ee622a │ 866 MiB │
│    │ e966a.asset/AssetData/Restore/022-10310-098.dmg                                                      │         │
│  3 │ /System/Library/AssetsV2/com_apple_MobileAsset_MacSoftwareUpdate/422e7e556d3f01a4bbf259b0f10f7ee622a │ 582 MiB │
│    │ e966a.asset/AssetData/Restore/BaseSystem.dmg                                                         │         │
│  4 │ /Users/egr/OneDrive/Фото/12131600_2294.MOV                                                           │ 481 MiB │
│  5 │ /System/Library/Speech/Voices/YelenaSiri.SpeechVoice/Contents/Resources/adat                         │ 414 MiB │
│  6 │ /System/Library/Speech/Voices/AlexCompact.SpeechVoice/Contents/Resources/PCMWave                     │ 306 MiB │
│  7 │ /usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/.git/objects/pack/pack-f507444fc7985877bdac2 │ 304 MiB │
│    │ 46d10cad99e7aca0cee.pack                                                                             │         │
│  8 │ /Users/egr/Library/Group Containers/3L68KQB4HG.group.com.readdle.smartemail/databases/cache.sqlite   │ 278 MiB │
│  9 │ /Users/egr/OneDrive/Изображения/Пленка/VID_20160608_212133.mp4                                       │ 239 MiB │
│ 10 │ /Users/egr/OneDrive/Изображения/Пленка/20140330_225521_Android.mp4                                   │ 222 MiB │
└────┴──────────────────────────────────────────────────────────────────────────────────────────────────────┴─────────┘

TOP 10 folders by size:

┌────┬──────────────────────────────────────────────────────────────────────────────────────────────────────┬───────┬───────┬─────────┬───────┐
│  # │ FOLDER                                                                                               │ FILES │ %     │ SIZE    │ %     │
├────┼──────────────────────────────────────────────────────────────────────────────────────────────────────┼───────┼───────┼─────────┼───────┤
│  1 │ /Users/egr/OneDrive/Изображения/Пленка                                                               │ 765   │ 0.08% │ 2.8 GiB │ 3.47% │
│  2 │ /Users/egr/OneDrive/Документы                                                                        │ 285   │ 0.03% │ 1.6 GiB │ 1.99% │
│  3 │ /System/Library/AssetsV2/com_apple_MobileAsset_MacSoftwareUpdate/422e7e556d3f01a4bbf259b0f10f7ee622a │ 6     │ 0.00% │ 1.4 GiB │ 1.76% │
│    │ e966a.asset/AssetData/Restore                                                                        │       │       │         │       │
│  4 │ /private/var/vm                                                                                      │ 1     │ 0.00% │ 1.0 GiB │ 1.24% │
│  5 │ /Users/egr/Library/Application Support/Microsoft/Teams/Service Worker/CacheStorage/2b5c392d2730c0910 │ 1,611 │ 0.17% │ 753 MiB │ 0.91% │
│    │ fd56433cc5e73e510d0f2b4/120aca5a-7991-4032-919f-02b360b9ed9f                                         │       │       │         │       │
│  6 │ /Library/Developer/CommandLineTools/usr/bin                                                          │ 84    │ 0.01% │ 738 MiB │ 0.89% │
│  7 │ /Applications/Adobe Premiere Pro 2020/Adobe Premiere Pro 2020.app/Contents/Frameworks                │ 70    │ 0.01% │ 629 MiB │ 0.76% │
│  8 │ /Applications/Adobe Media Encoder 2020/Adobe Media Encoder 2020.app/Contents/Frameworks              │ 69    │ 0.01% │ 624 MiB │ 0.76% │
│  9 │ /System/Library/Speech/Voices/YelenaSiri.SpeechVoice/Contents/Resources                              │ 10    │ 0.00% │ 550 MiB │ 0.67% │
│ 10 │ /Users/egr/OneDrive/Фото                                                                             │ 11    │ 0.00% │ 521 MiB │ 0.63% │
└────┴──────────────────────────────────────────────────────────────────────────────────────────────────────┴───────┴───────┴─────────┴───────┘

TOP 10 folders by count:

┌────┬──────────────────────────────────────────────────────────────────────────────────────────────────────┬────────┬───────┬─────────┬───────┐
│  # │ FOLDER                                                                                               │ FILES  │ %     │ SIZE    │ %     │
├────┼──────────────────────────────────────────────────────────────────────────────────────────────────────┼────────┼───────┼─────────┼───────┤
│  1 │ /Library/Developer/CommandLineTools/SDKs/MacOSX11.1.sdk/usr/share/man/man3                           │ 15,280 │ 1.61% │ 174 MiB │ 0.21% │
│  2 │ /Users/egr/Library/Containers/com.apple.iBooksX/Data/Library/Caches/WebKit/NetworkCache/Version 16/R │ 13,023 │ 1.37% │ 199 MiB │ 0.24% │
│    │ ecords/529573B272B0F7181050ED6AEAAD6113AFC26BF3/Resource                                             │        │       │         │       │
│  3 │ /Library/Developer/CommandLineTools/SDKs/MacOSX10.15.sdk/usr/share/man/man3                          │ 11,648 │ 1.23% │ 127 MiB │ 0.15% │
│  4 │ /Users/egr/Library/Application Support/Microsoft/Teams/Cache                                         │ 9,813  │ 1.03% │ 292 MiB │ 0.35% │
│  5 │ /Users/egr/Library/Containers/com.apple.iBooksX/Data/Library/Caches/WebKit/NetworkCache/Version 16/B │ 5,931  │ 0.63% │ 157 MiB │ 0.19% │
│    │ lobs                                                                                                 │        │       │         │       │
│  6 │ /usr/local/Homebrew/Library/Taps/homebrew/homebrew-core/Formula                                      │ 5,448  │ 0.57% │ 9.9 MiB │ 0.01% │
│  7 │ /Applications/GarageBand.app/Contents/Frameworks/MAResources.framework/Versions/A/Resources          │ 4,618  │ 0.49% │ 148 MiB │ 0.18% │
│  8 │ /usr/local/Homebrew/Library/Taps/homebrew/homebrew-cask/Casks                                        │ 3,791  │ 0.40% │ 2.2 MiB │ 0.00% │
│  9 │ /Applications/Adobe Premiere Pro 2020/Adobe Premiere Pro 2020.app/Contents/Frameworks/Frontend.frame │ 3,575  │ 0.38% │ 12 MiB  │ 0.01% │
│    │ work/Versions/A/Resources/png                                                                        │        │       │         │       │
│ 10 │ /Applications/Adobe Media Encoder 2020/Adobe Media Encoder 2020.app/Contents/Frameworks/AMEFrontend. │ 3,445  │ 0.36% │ 5.2 MiB │ 0.01% │
│    │ framework/Versions/A/Resources/png                                                                   │        │       │         │       │
└────┴──────────────────────────────────────────────────────────────────────────────────────────────────────┴────────┴───────┴─────────┴───────┘

Total files:            948,748 (81 GiB)
Total folders:          166,491
Total file extensions:  2,235

Read taken:	13.634901299s
```
