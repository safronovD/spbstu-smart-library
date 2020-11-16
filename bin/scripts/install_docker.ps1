$Packages = 'docker', 'docker-machine', 'virtualbox', 'docker-compose'
ForEach ($PackageName in $Packages)
{
    choco install $PackageName -y
}

