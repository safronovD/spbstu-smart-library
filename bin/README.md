## Usage

!! RUN AS ADMINISTRATOR

1) test_install.bat - install chocolatey, virtualbox, vagrant

There is some trouble with vagrant package - very long installation proccess
Close window after at least 5 minutes and try next step
If errors will be get - repeat this step

2) Restart your pc
3) test_start.bat - starting application on virtual machine

You can acces application on printed address
Go Menu -> Kibana -> Visualize -> Choose set -> Configure timeline

4) After working with application use Ctrl+C for stopping
5) test_clean.bat - cleaning virtual machine and uninstalling tools 
## Old version
!! Run as Administrator 
1) install_soft.bat - install chocolatey, virtualbox, docker-compose, docker and docker-machine
2) install_all.bat - install above software and setup ElasticSearch and Kibana
3) setup_ek.bat - deploy ElasticSearch and Kibana locally
4) cleanup_all.bat - delete all installed software.