docker-machine create -d virtualbox dev
docker-machine env --shell powershell dev | iex
Set-Location .\scripts\config
docker-compose up -d --build

Start-Sleep -Seconds 180
$str = docker-machine env dev | out-string -stream | findstr /r "[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*\.[0-9][0-9]*"
$ip = $str.split('/')[2].split(':')[0]
$url = 'http://' + $ip + ':5601/' 
Start-Process $url
