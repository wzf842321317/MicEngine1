@echo off
color 1f
 taskkill /f /im MicEngine.exe
 net stop  "MicEngine"
 ping 127.0.0.1 -n 3
 net start  "MicEngine"
pause