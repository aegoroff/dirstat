dirstat
=======

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/771c99225c8c4ad0a83115d2403f9fa0)](https://app.codacy.com/manual/egoroff/dirstat?utm_source=github.com&utm_medium=referral&utm_content=aegoroff/dirstat&utm_campaign=Badge_Grade_Dashboard) [![codecov](https://codecov.io/gh/aegoroff/dirstat/branch/master/graph/badge.svg)](https://codecov.io/gh/aegoroff/dirstat)

Small tool that shows selected folder or drive (on Windows) usage statistic.
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
  b           Show the first digit distribution of files size (benford law validation)
  e           Show file extensions statistic
  fi          Show information only about files
  fo          Show information only about folders
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
 1    Between 0 B and 100 KiB        1947384    89.52%    19 GiB     5.18%
 2    Between 100 KiB and 1.0 MiB    187479     8.62%     55 GiB     15.10%
 3    Between 1.0 MiB and 10 MiB     36246      1.67%     93 GiB     25.54%
 4    Between 10 MiB and 100 MiB     4451       0.20%     114 GiB    31.25%
 5    Between 100 MiB and 1.0 GiB    154        0.01%     34 GiB     9.23%
 6    Between 1.0 GiB and 10 GiB     6          0.00%     26 GiB     7.15%
 7    Between 10 GiB and 100 GiB     1          0.00%     24 GiB     6.73%
 8    Between 100 GiB and 1.0 TiB    0          0.00%     0 B        0.00%
 9    Between 1.0 TiB and 10 TiB     0          0.00%     0 B        0.00%
10    Between 10 TiB and 1.0 PiB     0          0.00%     0 B        0.00%

The first file's size digit distribution of non zero files (benford law):

Digit    Count     %         Benford ideal    %            Deviation
-----    -----     ------    -------------    ---------    ---------
1        669892    31.44%    641409           30.10%       4.44%
2        382144    17.93%    375043           17.60%       1.89%
3        268252    12.59%    266366           12.50%       0.71%
4        207326    9.73%     206700           9.70%        0.30%
5        160745    7.54%     168343           7.90%        -4.51%
6        131332    6.16%     142772           6.70%        -8.01%
7        108693    5.10%     123593           5.80%        -12.06%
8        113129    5.31%     108677           5.10%        4.10%
9        89416     4.20%     98022            4.60%        -8.78%

TOP 10 file extensions by size:

 #    Extension    Count     %         Size      %
--    ---------    -----     ------    ----      ------
 1    .dll         229160    10.53%    82 GiB    22.52%
 2    .vhdx        5         0.00%     37 GiB    10.06%
 3    .exe         18039     0.83%     21 GiB    5.71%
 4    .lib         16130     0.74%     20 GiB    5.57%
 5    .xml         99971     4.60%     14 GiB    3.93%
 6    .msi         2870      0.13%     12 GiB    3.37%
 7                 192647    8.86%     12 GiB    3.32%
 8    .cab         2496      0.11%     12 GiB    3.31%
 9    .pdb         14303     0.66%     11 GiB    2.97%
10    .jar         10122     0.47%     11 GiB    2.97%

TOP 10 file extensions by count:

 #    Extension    Count     %         Size       %
--    ---------    -----     ------    ----       ------
 1    .dll         229160    10.53%    82 GiB     22.52%
 2                 192647    8.86%     12 GiB     3.32%
 3    .js          167476    7.70%     3.3 GiB    0.91%
 4    .xml         99971     4.60%     14 GiB     3.93%
 5    .png         81947     3.77%     929 MiB    0.25%
 6    .py          78045     3.59%     816 MiB    0.22%
 7    .h           75250     3.46%     2.9 GiB    0.79%
 8    .cs          59269     2.72%     644 MiB    0.17%
 9    .go          54510     2.51%     785 MiB    0.21%
10    .html        51568     2.37%     804 MiB    0.22%

TOP 10 files by size:

 #    File                                                                                                                 Size
--    ------                                                                                                               ----
 1    c:\Users\egr\AppData\Local\Docker\wsl\data\ext4.vhdx                                                                 24 GiB
 2    c:\ProgramData\DockerDesktop\vm-data\DockerDesktop.vhdx                                                              9.5 GiB
 3    c:\Users\egr\java_error_in_rider.hprof                                                                               8.2 GiB
 4    c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0\flash.vhd                                       2.6 GiB
 5    c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.Ubuntu20.04onWindows_79rhkp1fndgsc\LocalState\ext4.vhdx    2.5 GiB
 6    c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.0\Emulation\Images\Flash.vhd                                  1.9 GiB
 7    c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.1\Emulation\Images\flash.vhd                                  1.5 GiB
 8    c:\Program Files\MongoDB\Server\4.0\data\collection-7--2807253782407563390.wt                                        962 MiB
 9    c:\Program Files\MongoDB\Server\4.0\data\collection-2-2535253104322579421.wt                                         902 MiB
10    c:\Windows\Installer\13b2694.msi                                                                                     887 MiB

TOP 10 folders by size:

 #    Folder                                                                                                     Files    %         Size       %
--    ------                                                                                                     -----    ------    ----       ------
 1    c:\Users\egr\AppData\Local\Docker\wsl\data                                                                 1        0.00%     24 GiB     6.73%
 2    c:\Windows\Installer                                                                                       2879     0.13%     14 GiB     3.71%
 3    c:\ProgramData\DockerDesktop\vm-data                                                                       1        0.00%     9.5 GiB    2.61%
 4    c:\Users\egr                                                                                               62       0.00%     8.3 GiB    2.27%
 5    c:\Program Files\MongoDB\Server\4.0\data                                                                   39       0.00%     3.4 GiB    0.92%
 6    c:\Windows\System32                                                                                        5082     0.23%     2.6 GiB    0.72%
 7    c:\Program Files (x86)\Windows Kits\10\Emulation\Mobile\10.0.10586.0                                       1        0.00%     2.6 GiB    0.70%
 8    c:\Users\egr\AppData\Local\Packages\CanonicalGroupLimited.Ubuntu20.04onWindows_79rhkp1fndgsc\LocalState    1        0.00%     2.5 GiB    0.68%
 9    c:\code\openssl\test                                                                                       1253     0.06%     2.1 GiB    0.58%
10    c:\Program Files (x86)\Microsoft SDKs\Windows Phone\v8.0\Emulation\Images                                  5        0.00%     2.0 GiB    0.54%

TOP 10 folders by count:

 #    Folder                                                                                             Files    %         Size       %
--    ------                                                                                             -----    ------    ----       ------
 1    c:\Windows\WinSxS\Manifests                                                                        23252    1.07%     21 MiB     0.01%
 2    c:\Windows\servicing\LCU\Package_for_RollupFix~31bf3856ad364e35~amd64~~19041.450.1.7               19865    0.91%     376 MiB    0.10%
 3    c:\Windows\servicing\LCU\Package_for_RollupFix~31bf3856ad364e35~amd64~~19041.388.1.7               13655    0.63%     331 MiB    0.09%
 4    c:\Users\egr\AppData\Local\Microsoft\Edge\User Data\Default\Code Cache\js                          12903    0.59%     281 MiB    0.08%
 5    c:\Users\egr\AppData\Local\Google\Chrome\User Data\Default\Local Storage                           6864     0.32%     84 MiB     0.02%
 6    c:\msys64\usr\share\man\man3                                                                       6739     0.31%     24 MiB     0.01%
 7    c:\Users\egr\AppData\Roaming\Microsoft\Crypto\RSA\S-1-5-21-2126183571-1528345743-823968152-1001    5144     0.24%     15 MiB     0.00%
 8    c:\Windows\System32                                                                                5082     0.23%     2.6 GiB    0.72%
 9    c:\Windows\servicing\Packages                                                                      4132     0.19%     59 MiB     0.02%
10    c:\msys64\usr\share\doc\openssl\html\man3                                                          3792     0.17%     36 MiB     0.01%

Total files:            2175304 (364 GiB)
Total folders:          459026
Total file extensions:  7622

Read taken:     8.5930367s

Alloc = 39 MiB  TotalAlloc = 3.5 GiB    Sys = 748 MiB   NumGC = 122     NumGoRoutines = 1
```