# coding: utf-8
Vagrant.configure("2") do |config|
  config.vm.box = "ubuntu/trusty64"

  config.vm.provider "virtualbox" do |vb|
    vb.gui = false
    vb.memory = "2048"
  end

  config.vm.define "release" do |release|
    release.vm.provision "shell", inline: <<-SHELL
cd ~
wget –quiet -c https://storage.googleapis.com/golang/go1.7.1.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.7.1.linux-amd64.tar.gz
echo 'export PATH=/usr/local/go/bin:$PATH' >> /home/vagrant/.bashrc
  SHELL
  end

  config.vm.define "built" do |built|
    built.vm.provision "shell", inline: <<-SHELL
cd ~
wget –quiet -c https://storage.googleapis.com/golang/go1.6.3.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.6.3.linux-amd64.tar.gz

apt-get update
apt-get install -y build-essential gcc git make

cd /home/vagrant
git clone https://go.googlesource.com/go
cd go
git checkout go1.7.1
cd src
GOROOT_BOOTSTRAP=/usr/local/go CGO_ENABLED=0 ./all.bash

echo 'export PATH=/home/vagrant/go/bin:$PATH' >> /home/vagrant/.bashrc
  SHELL
  end

  config.vm.define "docker" do |docker|
    docker.vm.provision "shell", inline: <<-SHELL
apt-get update
apt-get install -y apt-transport-https ca-certificates
apt-key adv --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys 58118E89F3A912897C070ADBF76221572C52609D
echo "deb https://apt.dockerproject.org/repo ubuntu-trusty main" > /etc/apt/sources.list.d/docker.list
apt-get update
apt-get install -y linux-image-extra-$(uname -r) linux-image-extra-virtual
apt-get install -y docker-engine
  SHELL
  end
  
end
