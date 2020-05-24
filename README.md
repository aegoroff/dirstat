dirstat
=======

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
  help        Help about any command
  version     Print the version number of dirstat

Flags:
  -h, --help          help for dirstat
  -p, --path string   REQUIRED. Directory path to show info.
  -r, --range ints    Output verbose files info for fileSizeRanges specified
  -v, --verbose       Be verbose

Use "dirstat [command] --help" for more information about a command.
```
Examples:
---------
Show statistic about all D drive (on Windows)
```
dirstat -p d:
```
Show statistic about specific folder d:\folder
```
dirstat -p d:\folder
```
Show statistic and additional verbose statistic about files in ranges Between 10 GiB and 100 GiB and Between 100 GiB and 1.0 TiB 
```
dirstat -p d:\ -v -r 7 -r 8
```
Output example:
---------------
```
Root: c:\

Total files stat:

File size                         Amount     %         Size       %
---------                         ------     ------    ----       ------
1. Between 0 B and 100 KiB        2312623    90.76%    20 GiB     6.00%
2. Between 100 KiB and 1.0 MiB    194939     7.65%     56 GiB     16.46%
3. Between 1.0 MiB and 10 MiB     36456      1.43%     93 GiB     27.14%
4. Between 10 MiB and 100 MiB     4477       0.18%     115 GiB    33.57%
5. Between 100 MiB and 1.0 GiB    157        0.01%     34 GiB     9.92%
6. Between 1.0 GiB and 10 GiB     6          0.00%     24 GiB     7.10%
7. Between 10 GiB and 100 GiB     0          0.00%     0 B        0.00%
8. Between 100 GiB and 1.0 TiB    0          0.00%     0 B        0.00%
9. Between 1.0 TiB and 10 TiB     0          0.00%     0 B        0.00%
10. Between 10 TiB and 100 TiB    0          0.00%     0 B        0.00%

TOP 10 file extensions by size:

Extension    Count     %         Size       %
---------    -----     ------    ----       ------
.dll         229492    9.01%     82 GiB     24.09%
.lib         16087     0.63%     20 GiB     5.95%
.exe         19251     0.76%     19 GiB     5.56%
             372370    14.61%    15 GiB     4.28%
.cab         2548      0.10%     14 GiB     4.13%
.xml         100923    3.96%     14 GiB     4.03%
.msi         2909      0.11%     12 GiB     3.66%
.ipch        383       0.02%     9.8 GiB    2.86%
.pdb         11493     0.45%     9.7 GiB    2.84%
.vhdx        1         0.00%     9.4 GiB    2.76%

TOP 10 file extensions by count:

Extension    Count     %         Size       %
---------    -----     ------    ----       ------
             372370    14.61%    15 GiB     4.28%
.dll         229492    9.01%     82 GiB     24.09%
.js          178578    7.01%     3.3 GiB    0.98%
.png         103099    4.05%     1.3 GiB    0.37%
.xml         100923    3.96%     14 GiB     4.03%
.py          78971     3.10%     833 MiB    0.24%
.h           77835     3.05%     2.9 GiB    0.85%
.cs          57984     2.28%     640 MiB    0.18%
.manifest    57351     2.25%     508 MiB    0.15%
.html        55709     2.19%     804 MiB    0.23%

TOP 10 files by size:

File                                                                                                      Size
------                                                                                                    ----
1. c:\ProgramData\DockerDesktop\vm-data\DockerDesktop.vhdx                                                9.4 GiB
2. c:\ProgramData\Microsoft\Search\Data\Applications\Windows\Windows.edb                                  4.9 GiB
3. c:\ProgramData\Microsoft\Windows\Hyper-V\Virtual Machines\2C2CA1B9-FF6C-4C0C-B5FE-4DC9FB59E1F2.VMRS    4.0 GiB
4. c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0\flash.vhd                         2.6 GiB
5. c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.0\Emulation\Images\Flash.vhd                    1.9 GiB
6. c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.1\Emulation\Images\flash.vhd                    1.5 GiB
7. c:\Program Files\MongoDB\Server\4.0\data\collection-2-2535253104322579421.wt                           892 MiB
8. c:\Windows\Installer\13b2694.msi                                                                       887 MiB
9. c:\Windows\Installer\9da85.msp                                                                         779 MiB
10. c:\code\Invoicing\.git\objects\pack\pack-ef3c358ee76889fc0180135924dfed992de72e63.pack                734 MiB

TOP 10 folders by size:

Folder                                                                          Files    %         Size       %
------                                                                          -----    ------    ----       ------
1. c:\Windows\Installer                                                         2869     0.11%     14 GiB     4.01%
2. c:\ProgramData\DockerDesktop\vm-data                                         1        0.00%     9.4 GiB    2.76%
3. c:\ProgramData\Microsoft\Search\Data\Applications\Windows                    12       0.00%     4.9 GiB    1.45%
4. c:\ProgramData\Microsoft\Windows\Hyper-V\Virtual Machines                    3        0.00%     4.0 GiB    1.17%
5. c:\Windows\System32                                                          5074     0.20%     2.7 GiB    0.78%
6. c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0         1        0.00%     2.6 GiB    0.75%
7. c:\code\openssl\test                                                         1253     0.05%     2.1 GiB    0.62%
8. c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.0\Emulation\Images    5        0.00%     2.0 GiB    0.58%
9. c:\Program Files\MongoDB\Server\4.0\data                                     35       0.00%     1.7 GiB    0.50%
10. c:\code\Invoicing\.git\objects\pack                                         46       0.00%     1.6 GiB    0.48%

Total files:            2548196 (342 GiB)
Total folders:          497637
Total file extensions:  10868

Read taken:    7.8638908s

Alloc = 159 MiB TotalAlloc = 3.7 GiB    Sys = 1.5 GiB   NumGC = 45
```