@echo off
color 1f
 taskkill /f /im micdog.exe
 net stop  "MicDog"
 ping 127.0.0.1 -n 3
 net start  "MicDog"
pause