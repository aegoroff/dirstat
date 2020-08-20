dirstat
=======

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/771c99225c8c4ad0a83115d2403f9fa0)](https://app.codacy.com/manual/egoroff/dirstat?utm_source=github.com&utm_medium=referral&utm_content=aegoroff/dirstat&utm_campaign=Badge_Grade_Dashboard) [![codecov](https://codecov.io/gh/aegoroff/dirstat/branch/master/graph/badge.svg)](https://codecov.io/gh/aegoroff/dirstat)

A small tool that shows selected folder or drive (on Windows) usage statistic.
The tool shows file and folder statistic like the number of files or folders, top 10 files and folders by size,
file statistics by extension and by file size range.

Syntax:
------
```
 A small tool that shows selected folder or drive (on Windows) usage statistic

Usage:
  dirstat [flags]
  dirstat [command]

Available Commands:
  a           Show all information about folder/volume
  fi          Show information about files within folder on volume only
  fo          Show information about folders within folder on volume only
  help        Help about any command
  version     Print the version number of dirstat

Flags:
  -h, --help         help for dirstat
  -m, --memory       Show memory statistic after run
  -o, --removeroot   Remove root part from full path i.e. output relative paths
  -t, --top int      The number of lines in top statistics. (default 10)

Use "dirstat [command] --help" for more information about a command.
```
Examples:
---------
Show statistic about all D drive (on Windows)
```
dirstat a -p d:
```
Show all statistic about specific folder d:\folder
```
dirstat a -p d:\folder
```
Show all statistic and additional verbose statistic about files in ranges Between 10 GiB and 100 GiB and Between 100 GiB and 1.0 TiB 
```
dirstat a -p d:\ -r 7 -r 8
```
or
```
dirstat a -p d:\ -r 7,8
```
The second form is equivalent

Output example:
---------------
```
Root: c:\

Total files stat:

 #    File size                      Amount     %         Size       %
--    ---------                      ------     ------    ----       ------
 1    Between 0 B and 100 KiB        1889276    89.45%    18 GiB     5.30%
 2    Between 100 KiB and 1.0 MiB    183373     8.68%     54 GiB     15.44%
 3    Between 1.0 MiB and 10 MiB     35261      1.67%     91 GiB     26.02%
 4    Between 10 MiB and 100 MiB     4366       0.21%     112 GiB    32.05%
 5    Between 100 MiB and 1.0 GiB    145        0.01%     32 GiB     9.19%
 6    Between 1.0 GiB and 10 GiB     5          0.00%     18 GiB     5.13%
 7    Between 10 GiB and 100 GiB     1          0.00%     24 GiB     7.04%
 8    Between 100 GiB and 1.0 TiB    0          0.00%     0 B        0.00%
 9    Between 1.0 TiB and 10 TiB     0          0.00%     0 B        0.00%
10    Between 10 TiB and 1.0 PiB     0          0.00%     0 B        0.00%

TOP 10 file extensions by size:

 #    Extension    Count     %         Size      %
--    ---------    -----     ------    ----      ------
 1    .dll         223227    10.57%    80 GiB    23.02%
 2    .vhdx        5         0.00%     37 GiB    10.52%
 3    .exe         17275     0.82%     20 GiB    5.84%
 4    .lib         16111     0.76%     20 GiB    5.82%
 5    .xml         97838     4.63%     14 GiB    4.04%
 6    .msi         2866      0.14%     12 GiB    3.52%
 7    .cab         2494      0.12%     12 GiB    3.44%
 8                 184215    8.72%     11 GiB    3.28%
 9    .pdb         13551     0.64%     10 GiB    3.02%
10    .ipch        404       0.02%     10 GiB    2.99%

TOP 10 file extensions by count:

 #    Extension    Count     %         Size       %
--    ---------    -----     ------    ----       ------
 1    .dll         223227    10.57%    80 GiB     23.02%
 2                 184215    8.72%     11 GiB     3.28%
 3    .js          168308    7.97%     3.3 GiB    0.95%
 4    .xml         97838     4.63%     14 GiB     4.04%
 5    .py          76695     3.63%     809 MiB    0.23%
 6    .png         76384     3.62%     919 MiB    0.26%
 7    .h           75181     3.56%     2.9 GiB    0.83%
 8    .cs          59030     2.79%     643 MiB    0.18%
 9    .html        51372     2.43%     800 MiB    0.22%
10    .pyc         50191     2.38%     524 MiB    0.15%

TOP 10 files by size:

 #    File                                                                                                                 Size
--    ------                                                                                                               ----
 1    c:\Users\egr\AppData\Local\Docker\wsl\data\ext4.vhdx                                                                 24 GiB
 2    c:\ProgramData\DockerDesktop\vm-data\DockerDesktop.vhdx                                                              9.5 GiB
 3    c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0\flash.vhd                                       2.6 GiB
 4    c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.Ubuntu20.04onWindows_79rhkp1fndgsc\LocalState\ext4.vhdx    2.5 GiB
 5    c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.0\Emulation\Images\Flash.vhd                                  1.9 GiB
 6    c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.1\Emulation\Images\flash.vhd                                  1.5 GiB
 7    c:\Program Files\MongoDB\Server\4.0\data\collection-7--2807253782407563390.wt                                        962 MiB
 8    c:\Program Files\MongoDB\Server\4.0\data\collection-2-2535253104322579421.wt                                         902 MiB
 9    c:\Windows\Installer\13b2694.msi                                                                                     887 MiB
10    c:\Windows\Installer\9da85.msp                                                                                       779 MiB

TOP 10 folders by size:

 #    Folder                                                                                                       Files    %         Size       %
--    ------                                                                                                       -----    ------    ----       ------
 1    c:\Users\egr\AppData\Local\Docker\wsl\data                                                                   1        0.00%     24 GiB     7.04%
 2    c:\Windows\Installer                                                                                         2882     0.14%     14 GiB     4.03%
 3    c:\ProgramData\DockerDesktop\vm-data                                                                         1        0.00%     9.5 GiB    2.72%
 4    c:\Program Files\MongoDB\Server\4.0\data                                                                     39       0.00%     3.3 GiB    0.96%
 5    c:\Windows\System32                                                                                          5080     0.24%     2.6 GiB    0.76%
 6    c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0                                         1        0.00%     2.6 GiB    0.73%
 7    c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.Ubuntu20.04onWindows_79rhkp1fndgsc\LocalState      1        0.00%     2.5 GiB    0.71%
 8    c:\code\openssl\test                                                                                         1253     0.06%     2.1 GiB    0.61%
 9    c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.0\Emulation\Images                                    5        0.00%     2.0 GiB    0.57%
10    c:\ProgramData\Package Cache\{EA923538-8370-4294-A5CC-F6130FAAD89D}v10.1.10586.11\Redistributable\1.0.0.0    51       0.00%     1.5 GiB    0.43%

TOP 10 folders by count:

 #    Folder                                                                                             Files    %         Size       %
--    ------                                                                                             -----    ------    ----       ------
 1    c:\Windows\WinSxS\Manifests                                                                        22722    1.08%     20 MiB     0.01%
 2    c:\Windows\servicing\LCU\Package_for_RollupFix~31bf3856ad364e35~amd64~~19041.388.1.7               13655    0.65%     331 MiB    0.09%
 3    c:\Windows\servicing\LCU\Package_for_RollupFix~31bf3856ad364e35~amd64~~19041.329.1.7               11665    0.55%     298 MiB    0.08%
 4    c:\Users\egr\AppData\Local\Microsoft\Edge\User Data\Default\Code Cache\js                          11472    0.54%     285 MiB    0.08%
 5    c:\Users\egr\AppData\Local\Google\Chrome\User Data\Default\Local Storage                           6864     0.32%     84 MiB     0.02%
 6    c:\msys64\usr\share\man\man3                                                                       6739     0.32%     24 MiB     0.01%
 7    c:\Users\egr\AppData\Roaming\Microsoft\Crypto\RSA\S-1-5-21-2126183571-1528345743-823968152-1001    5144     0.24%     15 MiB     0.00%
 8    c:\Windows\System32                                                                                5080     0.24%     2.6 GiB    0.76%
 9    c:\Windows\servicing\Packages                                                                      3838     0.18%     56 MiB     0.02%
10    c:\msys64\usr\share\doc\openssl\html\man3                                                          3792     0.18%     36 MiB     0.01%

Total files:            2112029 (348 GiB)
Total folders:          445144
Total file extensions:  7579

Read taken:	8.4079962s

Alloc = 44 MiB	TotalAlloc = 3.4 GiB	Sys = 714 MiB	NumGC = 117	NumGoRoutines = 1
```