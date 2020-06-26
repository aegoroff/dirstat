dirstat
=======

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/771c99225c8c4ad0a83115d2403f9fa0)](https://app.codacy.com/manual/egoroff/dirstat?utm_source=github.com&utm_medium=referral&utm_content=aegoroff/dirstat&utm_campaign=Badge_Grade_Dashboard)

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

File size                          Amount     %         Size       %
---------                          ------     ------    ----       ------
 1. Between 0 B and 100 KiB        1898300    89.50%    18 GiB     5.50%
 2. Between 100 KiB and 1.0 MiB    182968     8.63%     54 GiB     16.13%
 3. Between 1.0 MiB and 10 MiB     35430      1.67%     91 GiB     27.36%
 4. Between 10 MiB and 100 MiB     4432       0.21%     113 GiB    33.84%
 5. Between 100 MiB and 1.0 GiB    151        0.01%     32 GiB     9.61%
 6. Between 1.0 GiB and 10 GiB     7          0.00%     26 GiB     7.72%
 7. Between 10 GiB and 100 GiB     0          0.00%     0 B        0.00%
 8. Between 100 GiB and 1.0 TiB    0          0.00%     0 B        0.00%
 9. Between 1.0 TiB and 10 TiB     0          0.00%     0 B        0.00%
10. Between 10 TiB and 1.0 PiB     0          0.00%     0 B        0.00%

TOP 10 file extensions by size:

Extension    Count     %         Size       %
---------    -----     ------    ----       ------
.dll         219033    10.33%    82 GiB     24.50%
.lib         16109     0.76%     20 GiB     6.11%
.exe         17194     0.81%     20 GiB     6.00%
.vhdx        5         0.00%     16 GiB     4.70%
.cab         2543      0.12%     14 GiB     4.19%
.xml         98048     4.62%     14 GiB     4.16%
.msi         2911      0.14%     12 GiB     3.72%
             183501    8.65%     12 GiB     3.53%
.pdb         12072     0.57%     10 GiB     2.99%
.ipch        390       0.02%     9.9 GiB    2.97%

TOP 10 file extensions by count:

Extension    Count     %         Size       %
---------    -----     ------    ----       ------
.dll         219033    10.33%    82 GiB     24.50%
             183501    8.65%     12 GiB     3.53%
.js          173216    8.17%     3.3 GiB    0.99%
.xml         98048     4.62%     14 GiB     4.16%
.png         88031     4.15%     992 MiB    0.29%
.py          76470     3.61%     805 MiB    0.24%
.h           74113     3.49%     2.9 GiB    0.86%
.cs          59033     2.78%     642 MiB    0.19%
.html        51327     2.42%     801 MiB    0.23%
.pyc         50191     2.37%     524 MiB    0.15%

TOP 10 files by size:

File                                                                                                                Size
------                                                                                                              ----
 1. c:\ProgramData\DockerDesktop\vm-data\DockerDesktop.vhdx                                                         9.5 GiB
 2. c:\ProgramData\Microsoft\Search\Data\Applications\Windows\Windows.edb                                           4.3 GiB
 3. c:\Users\egr\AppData\Local\Docker\wsl\data\ext4.vhdx                                                            3.6 GiB
 4. c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0\flash.vhd                                  2.6 GiB
 5. c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.UbuntuonWindows_79rhkp1fndgsc\LocalState\ext4.vhdx    2.4 GiB
 6. c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.0\Emulation\Images\Flash.vhd                             1.9 GiB
 7. c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.1\Emulation\Images\flash.vhd                             1.5 GiB
 8. c:\Program Files\MongoDB\Server\4.0\data\collection-2-2535253104322579421.wt                                    891 MiB
 9. c:\Windows\Installer\13b2694.msi                                                                                887 MiB
10. c:\Windows\Installer\9da85.msp                                                                                  779 MiB

TOP 10 folders by size:

Folder                                                                                                    Files    %         Size       %
------                                                                                                    -----    ------    ----       ------
 1. c:\Windows\Installer                                                                                  2872     0.14%     13 GiB     4.02%
 2. c:\ProgramData\DockerDesktop\vm-data                                                                  1        0.00%     9.5 GiB    2.84%
 3. c:\ProgramData\Microsoft\Search\Data\Applications\Windows                                             10       0.00%     4.3 GiB    1.30%
 4. c:\Users\egr\AppData\Local\Docker\wsl\data                                                            1        0.00%     3.6 GiB    1.07%
 5. c:\Windows\System32                                                                                   5079     0.24%     2.7 GiB    0.81%
 6. c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0                                  1        0.00%     2.6 GiB    0.77%
 7. c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.UbuntuonWindows_79rhkp1fndgsc\LocalState    2        0.00%     2.4 GiB    0.73%
 8. c:\code\openssl\test                                                                                  1253     0.06%     2.1 GiB    0.63%
 9. c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.0\Emulation\Images                             5        0.00%     2.0 GiB    0.59%
10. c:\Program Files\MongoDB\Server\4.0\data                                                              35       0.00%     1.7 GiB    0.51%

TOP 10 folders by count:

Folder                                                                                                 Files    %         Size       %
------                                                                                                 -----    ------    ----       ------
 1. c:\Windows\WinSxS\Manifests                                                                        24056    1.13%     23 MiB     0.01%
 2. c:\Windows\servicing\LCU\Package_for_RollupFix~31bf3856ad364e35~amd64~~19041.329.1.7               11665    0.55%     298 MiB    0.09%
 3. c:\Users\egr\AppData\Local\Google\Chrome\User Data\Default\Local Storage                           6864     0.32%     84 MiB     0.02%
 4. c:\msys64\usr\share\man\man3                                                                       6739     0.32%     24 MiB     0.01%
 5. c:\Users\egr\AppData\Local\Microsoft\Edge\User Data\Default\Code Cache\js                          5383     0.25%     294 MiB    0.09%
 6. c:\Users\egr\AppData\Roaming\Microsoft\Crypto\RSA\S-1-5-21-2126183571-1528345743-823968152-1001    5144     0.24%     15 MiB     0.00%
 7. c:\Windows\System32                                                                                5079     0.24%     2.7 GiB    0.81%
 8. c:\Windows\servicing\Packages                                                                      4410     0.21%     74 MiB     0.02%
 9. c:\msys64\usr\share\doc\openssl\html\man3                                                          3792     0.18%     36 MiB     0.01%
10. c:\Program Files\Adobe\Adobe Premiere Rush\PNG                                                     3593     0.17%     18 MiB     0.01%

Total files:            2120914 (333 GiB)
Total folders:          437235
Total file extensions:  9804

Read taken:    8.0340003s

Alloc = 142 MiB TotalAlloc = 3.4 GiB    Sys = 732 MiB   NumGC = 87
```