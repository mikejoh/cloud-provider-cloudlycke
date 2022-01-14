IMAGE_NAME = "ubuntu/bionic64"
CLUSTERS = 2
NODES = 1
NODE_CPU = 2
NODE_MEMORY = 1024
MASTER_CPU = 1
MASTER_MEMORY = 512

CLOUD_PROVIDER = "cloudlycke"

K8S_VERSION = "1.18.2"

Vagrant.configure("2") do |config|
    config.ssh.insert_key = true

    (1..CLUSTERS).each do |c|
        config.vm.define "master-c#{c}-1" do |master|
            config.vm.provider "virtualbox" do |master_vm|
              master_vm.cpus = MASTER_CPU
              master_vm.memory = MASTER_MEMORY
              master_vm.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ] # Disable console logging
            end
            master.vm.box = IMAGE_NAME
            master.vm.network "private_network", ip: "192.168.#{c}0.10"
            master.vm.hostname = "master-c#{c}-1"
            master.vm.provision "ansible" do |ansible|
                ansible.playbook = "vagrant/ansible/master.yaml"
                ansible.extra_vars = {
                    node_ip: "192.168.#{c}0.10",
                    node_name: "master-c#{c}-1",
                    node_id: "m-c#{c}-1",
                    cluster: "#{c}",
                    cloud_provider: CLOUD_PROVIDER,
                    k8s_version: K8S_VERSION
                }
            end
        end

        (1..NODES).each do |i|
            config.vm.define "node-c#{c}-#{i}" do |node|
                config.vm.provider "virtualbox" do |node_vm|
                  node_vm.cpus = NODE_CPU
                  node_vm.memory = NODE_MEMORY
                  node_vm.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ] # Disable console logging
                end
                node.vm.box = IMAGE_NAME
                node.vm.network "private_network", ip: "192.168.#{c}0.#{i + 10}"
                node.vm.hostname = "node-c#{c}-#{i}"
                node.vm.provision "ansible" do |ansible|
                    ansible.playbook = "vagrant/ansible/node.yaml"
                    ansible.extra_vars = {
                        node_ip: "192.168.#{c}0.#{i + 10}",
                        node_name: "node-c#{c}-#{i}",
                        node_id: "n-c#{c}-#{i}",
                        cluster: "#{c}",
                        cloud_provider: CLOUD_PROVIDER,
                        k8s_version: K8S_VERSION
                    }
                end
            end
        end
    end
end