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
  a           Show all information about folder/volume
  fi          Show information about files within folder on volume only
  fo          Show information about folders within folder on volume only
  help        Help about any command
  version     Print the version number of dirstat

Flags:
  -h, --help      help for dirstat
  -m, --memory    Show memory statistic after run
  -t, --top int   The number of lines in top statistics. (default 10)

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
dirstat a -p d:\ -v -r 7 -r 8
```
Output example:
---------------
```
Root: c:\

Total files stat:

File size                         Amount     %         Size       %
---------                         ------     ------    ----       ------
1. Between 0 B and 100 KiB        1867470    89.45%    18 GiB     5.48%
2. Between 100 KiB and 1.0 MiB    180829     8.66%     53 GiB     16.13%
3. Between 1.0 MiB and 10 MiB     35258      1.69%     90 GiB     27.45%
4. Between 10 MiB and 100 MiB     4404       0.21%     112 GiB    34.04%
5. Between 100 MiB and 1.0 GiB    144        0.01%     32 GiB     9.68%
6. Between 1.0 GiB and 10 GiB     7          0.00%     24 GiB     7.39%
7. Between 10 GiB and 100 GiB     0          0.00%     0 B        0.00%
8. Between 100 GiB and 1.0 TiB    0          0.00%     0 B        0.00%
9. Between 1.0 TiB and 10 TiB     0          0.00%     0 B        0.00%
10. Between 10 TiB and 1.0 PiB    0          0.00%     0 B        0.00%

TOP 10 file extensions by size:

Extension    Count     %         Size       %
---------    -----     ------    ----       ------
.dll         214160    10.26%    80 GiB     24.41%
.lib         16110     0.77%     20 GiB     6.18%
.exe         16596     0.79%     20 GiB     5.92%
.vhdx        6         0.00%     15 GiB     4.63%
.xml         97561     4.67%     14 GiB     4.20%
.cab         2329      0.11%     14 GiB     4.19%
.msi         2909      0.14%     12 GiB     3.75%
             178372    8.54%     12 GiB     3.61%
.pdb         11895     0.57%     9.9 GiB    3.01%
.ipch        384       0.02%     9.8 GiB    2.97%

TOP 10 file extensions by count:

Extension    Count     %         Size       %
---------    -----     ------    ----       ------
.dll         214160    10.26%    80 GiB     24.41%
             178372    8.54%     12 GiB     3.61%
.js          172918    8.28%     3.3 GiB    1.01%
.xml         97561     4.67%     14 GiB     4.20%
.png         85731     4.11%     991 MiB    0.29%
.py          76470     3.66%     805 MiB    0.24%
.h           74107     3.55%     2.9 GiB    0.87%
.cs          58682     2.81%     641 MiB    0.19%
.html        51073     2.45%     802 MiB    0.24%
.pyc         50191     2.40%     524 MiB    0.16%

TOP 10 files by size:

File                                                                                                               Size
------                                                                                                             ----
1. c:\ProgramData\DockerDesktop\vm-data\DockerDesktop.vhdx                                                         9.5 GiB
2. c:\ProgramData\Microsoft\Search\Data\Applications\Windows\Windows.edb                                           4.3 GiB
3. c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0\flash.vhd                                  2.6 GiB
4. c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.UbuntuonWindows_79rhkp1fndgsc\LocalState\ext4.vhdx    2.4 GiB
5. c:\Users\egr\AppData\Local\Docker\wsl\data\ext4.vhdx                                                            2.3 GiB
6. c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.0\Emulation\Images\Flash.vhd                             1.9 GiB
7. c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.1\Emulation\Images\flash.vhd                             1.5 GiB
8. c:\Users\egr\AppData\Local\Packages\46932SUSE.openSUSELeap42.2_022rs5jcyhyac\LocalState\ext4.vhdx               971 MiB
9. c:\Program Files\MongoDB\Server\4.0\data\collection-2-2535253104322579421.wt                                    892 MiB
10. c:\Windows\Installer\13b2694.msi                                                                               887 MiB

TOP 10 folders by size:

Folder                                                                                                   Files    %         Size       %
------                                                                                                   -----    ------    ----       ------
1. c:\Windows\Installer                                                                                  2867     0.14%     13 GiB     4.06%
2. c:\ProgramData\DockerDesktop\vm-data                                                                  1        0.00%     9.5 GiB    2.87%
3. c:\ProgramData\Microsoft\Search\Data\Applications\Windows                                             10       0.00%     4.3 GiB    1.32%
4. c:\Windows\System32                                                                                   5078     0.24%     2.7 GiB    0.81%
5. c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0                                  1        0.00%     2.6 GiB    0.77%
6. c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.UbuntuonWindows_79rhkp1fndgsc\LocalState    2        0.00%     2.4 GiB    0.72%
7. c:\Users\egr\AppData\Local\Docker\wsl\data                                                            1        0.00%     2.3 GiB    0.69%
8. c:\code\openssl\test                                                                                  1253     0.06%     2.1 GiB    0.64%
9. c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.0\Emulation\Images                             5        0.00%     2.0 GiB    0.60%
10. c:\Program Files\MongoDB\Server\4.0\data                                                             35       0.00%     1.7 GiB    0.52%

TOP 10 folders by count:

Folder                                                                                                Files    %         Size       %
------                                                                                                -----    ------    ----       ------
1. c:\Windows\WinSxS\Manifests                                                                        23328    1.12%     21 MiB     0.01%
2. c:\Users\egr\AppData\Local\Google\Chrome\User Data\Default\Local Storage                           6864     0.33%     84 MiB     0.02%
3. c:\msys64\usr\share\man\man3                                                                       6739     0.32%     24 MiB     0.01%
4. c:\Users\egr\AppData\Roaming\Microsoft\Crypto\RSA\S-1-5-21-2126183571-1528345743-823968152-1001    5144     0.25%     15 MiB     0.00%
5. c:\Windows\System32                                                                                5078     0.24%     2.7 GiB    0.81%
6. c:\Windows\servicing\Packages                                                                      3908     0.19%     57 MiB     0.02%
7. c:\msys64\usr\share\doc\openssl\html\man3                                                          3792     0.18%     36 MiB     0.01%
8. c:\Program Files\Adobe\Adobe Premiere Rush\PNG                                                     3593     0.17%     18 MiB     0.01%
9. c:\Program Files\Adobe\Adobe Premiere Pro 2020\PNG                                                 3593     0.17%     18 MiB     0.01%
10. c:\Users\egr\AppData\Roaming\Adobe\Common\Media Cache Files                                       3372     0.16%     98 MiB     0.03%

Total files:            2087740 (330 GiB)
Total folders:          427517
Total file extensions:  9524

Read taken:    8.0809999s

Alloc = 181 MiB TotalAlloc = 3.3 GiB    Sys = 777 MiB   NumGC = 88
```