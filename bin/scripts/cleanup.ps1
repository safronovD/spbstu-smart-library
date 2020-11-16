docker-machine kill dev
$dir = 'C:\Users\' + $env:UserName + '\.docker\machine\machines\dev'
Remove-Item -path $dir -recurse 
$Packages = 'docker', 'docker-machine', 'virtualbox', 'docker-compose'
ForEach ($PackageName in $Packages)
{
    choco uninstall $PackageName -y
}

Remove-Item -path 'C:\ProgramData\chocolatey' -recurse -force

$path = [System.Environment]::GetEnvironmentVariable(
    'PATH',
    'Machine'
)

$path = ($path.Split(';') | Where-Object { $_ -ne 'C:\ProgramData\chocolatey\bin' }) -join ';'

[System.Environment]::SetEnvironmentVariable(
    'PATH',
    $path,
    'Machine'
)
