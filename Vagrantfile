# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure("2") do |config|
  config.vm.box = "bento/centos-7.8"
  config.vm.synced_folder '.', '/vagrant'

  config.vm.network 'forwarded_port',
    guest: 5601, host: 5601
  config.vm.network 'forwarded_port',
    guest: 9200, host: 9200

  config.vm.provision "Boostrap", type: "shell" do |sh|
    sh.inline = "echo Boostrap process #################"
    sh.inline = "sudo yum install -y docker;
                 sudo systemctl start docker;
                 sudo yum install -y python3"
  end

  config.vm.provision "Configurating compose", type: "shell" do |sh|
    sh.inline = "sudo pip3 install docker-compose;
                 sudo chmod +x /usr/local/bin/docker-compose;
                 sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose"
  end

  config.vm.provision "Starting application", type: "shell" do |sh|
    sh.inline = "cd /vagrant/bin/scripts/config;
                 sudo docker login -u maxvale -p c1dd4464542bf943885cb09bbbb63155a8632576 docker.pkg.github.com;
                 sudo docker-compose up"
  end

  config.vm.provider 'virtualbox' do |vb|
    vb.name = "Smart library"
    vb.memory = 2048
  end
  
end
