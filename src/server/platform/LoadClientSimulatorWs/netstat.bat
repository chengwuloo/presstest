@echo off
for /l %%i in (1,1,1000) do (
netstat -aon |find /V /C ""
netstat -aon |findstr "192.168.2.211" |find /V /C ""
echo.
ping /n 3 127.1>nul
)
pause