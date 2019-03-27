# escape=`
FROM plugins/base:windows-1809

LABEL maintainer="Drone.IO Community <drone-dev@googlegroups.com>" `
  org.label-schema.name="Drone Git" `
  org.label-schema.vendor="Drone.IO Community" `
  org.label-schema.schema-version="1.0"

RUN Invoke-WebRequest 'https://github.com/git-for-windows/git/releases/download/v2.12.2.windows.2/MinGit-2.12.2.2-64-bit.zip' -OutFile 'git.zip'; `
  Expand-Archive -Path git.zip -DestinationPath c:\git\ -Force; `
  $env:PATH = 'c:\git\cmd;c:\git\mingw64\bin;c:\git\usr\bin;{0}' -f $env:PATH; `
  Set-ItemProperty -Path 'HKLM:\SYSTEM\CurrentControlSet\Control\Session Manager\Environment\' -Name Path -Value $env:PATH; `
  Remove-Item -Path git.zip

ADD release/windows/amd64/drone-git.exe C:/bin/drone-git.exe
ENTRYPOINT [ "C:\\bin\\drone-git.exe" ]
