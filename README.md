# CMPE281 - Personal Project - Jay Parekh

## Project Details

Select one CP and one AP NoSQL database.

### Configuration

1. Set up your cluster as AWS EC2 Instances.

2. Set up the Experiments (i.e. Test Cases) to answer the following questions:

### Questions

1. How does the system function during normal mode (i.e. no partition)

2. What happens to the master node during a partition?

3. Can stale data be read from a slave node during a partition?

4. What happens to the system during partition recovery?

## **Journal**

### Submission Link

<https://youtu.be/K5uUphP7mLs>

### **MongoDB**

#### **General Flow for MongoDB:**

* Launch Linux EC2 instances.
* Install MongoDB on instances
* Connect them and initiate a replicaset.
* Create a partition and analyze the CAP properties.

#### **MongoDB Progress**

* [x] Create MongoDB Cluster
* [x] Test MongoDB Cluster
* [x] Create MongoDB Shards
* [x] Create GO API
* [x] Create Video

#### Create MongoDB Cluster

Here, we'll make a MongoDB cluster in AWS. Cluster size will be of 5 nodes. Total of 5 EC2 instances will be used and each instance will work as a node of the cluster.

***Reference**: <https://github.com/paulnguyen/cmpe281/blob/master/labs/lab4/aws-mongodb-replica-set.md>*

1. Create Jumpbox

    * AMI: Amazon Linux AMI 2018.03.0 (HVM), SSD Volume Type
    * Instance Type: t2.micro
    * Network: CMPE281
    * Subnet: Public Subnet
    * Auto-assign Public IP: Enable
    * Tag: jumpbox
    * Security Group: **jumpbox**
      * Ports: 22, 80
    * Keypair: cmpe281-us-west-1.pem

    **Note:** Since all the instances will be in private subnet, jumpbox is needed to access them.

1. Creating an EC2 instance for mongo
    * AMI: Ubuntu Server 16.04 LTS (HVM), SSD Volume Type
    * Instance Type: t2.micro
    * Network: CMPE281
    * Subnet: Private Subnet
    * Auto-assign Public IP: Disable
    * Tag: **mongo-primary**
    * Security Group: mongo
      * Ports: 22, 27017
    * Keypair: cmpe281-us-west-1.pem

    **Note:** MongoDB works on port 27017, we'll be communication this instance through port 27017.

1. Connecting to **mongo-primary**

    * Upload key to **jumpbox**

        ```bash
        scp -i cmpe281-us-west-1.pem cmpe281-us-west-1.pem ec2-user@ec2-13-56-16-49.us-west-1.compute.amazonaws.com:
        ```
    * Connect to **jumpbox**

        ```bash
        chmod 400 cmpe281-us-west-1.pem
        ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-13-56-161-186.us-west-1.compute.amazonaws.com
        ```
    * Connect to **mongo-primary**

        ```bash
        chmod 400 cmpe281-us-west-1.pem
        ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.67
        ```

1. Install MongoDB

    **Note:** Start the NAT-gateway instance of the VPC with Elastic IP in order to provide internet access to private subnet instances.

    1. Import the MongoDB repository

        * Import the public key used by the package management system.

            ```bash
            sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 9DA31620334BD75D9DCB49F368818C72E52529D4
            ```
        * Create a source list file for MongoDB

            ```bash
            echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/4.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb.list
            ```
        * Update the local package repository

            ```bash
            sudo apt update
            ```

    1. Install the MongoDB packages

        ```bash
        sudo apt install mongodb-org
        ```

    1. Launch MongoDB as a service
        * Enable MongoDB on startup

            ```bash
            sudo systemctl enable mongod
            ```

        * Start MongoDB service

            ```bash
            sudo systemctl start mongod 
            ```

        * Stop MongoDB service

            ```bash
            sudo systemctl stop mongod
            ```

        * Restart MongoDB service

            ```bash
            sudo systemctl restart mongod
            ```

        * Check status for MongoDB service

            ```bash
            sudo systemctl status mongod
            ```

1. Create MongoDB KeyFile

    ```bash
    openssl rand -base64 741 > keyFile
    sudo mkdir -p /opt/mongodb
    sudo cp keyFile /opt/mongodb
    sudo chown mongodb:mongodb /opt/mongodb/keyFile
    sudo chmod 0600 /opt/mongodb/keyFile
    ```

1. Config mongod.config

    Open mongod.config in edit mode

    ```bash
    sudo vi /etc/mongod.conf
    ```

    1. Set bindIp

        ```bash
        bindIp: 0.0.0.0
        ```

    1. Set keyfile as security

        ```bash
        security:
            keyFile: /opt/mongodb/keyFile
        ```

    1. Set replica set name

        ```bash
        replication:
            replSetName: cmpe281
        ```

1. Create mongod.service

    * Open file in edit mode

        ```bash
        sudo vi /etc/systemd/system/mongod.service
        ```

    * File Content

        ```bash
        [Unit]
            Description=High-performance, schema-free document-oriented database
            After=network.target

        [Service]
            User=mongodb
            ExecStart=/usr/bin/mongod --quiet --config /etc/mongod.conf

        [Install]
            WantedBy=multi-user.target
        ```

    * Enable Mongo Service

        ```bash
        sudo systemctl enable mongod.service
        ```

    * Restart MongoDB to apply our changes

        ```bash
        sudo service mongod restart
        ```

    * Check MongoDB status

        ```bash
        sudo service mongod status
        ```

1. Create Image of **mongo-primary**
    * Image Name: mongo
    * Image Description: mongo 4.0.4, ubuntu 16.04, replicaset=cmpe281

1. Launch Secondary Instances
    * AMI: mongo
    * Instance Type: t2.micro
    * Number of Instances: 5
    * Network: CMPE281
    * Subnet: Private Subnet
    * Auto-assign Public IP: Disable
    * Security Group: mongo
        * Ports: 22,27017
    * Key: cmpe281-us-west-2.pem
    * *Give them names mongo-secondary-1, mongo-secondary-2, mongo-secondary-3, mongo-secondary-4, mongo-secondary-5 for better understanding*

1. Information of Instances

    |Instance|IP|SSH|
    |--------|--|---|
    |jumpbox|13.56.161.186|ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-13-56-161-186.us-west-1.compute.amazonaws.com|
    |mongo-primary|10.0.1.115|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.115|
    |mongo-secondary-1|10.0.1.165|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.165|
    |mongo-secondary-2|10.0.1.175|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.175|
    |mongo-secondary-3|10.0.1.107|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.107|
    |mongo-secondary-4|10.0.1.211|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.211|

1. Changing the hostname of **jumpbox** for better understanding
    * Update the /etc/sysconfig/network file

        ```bash
        sudo vim /etc/sysconfig/network
            HOSTNAME=jumpbox
            NETWORKING=yes
        ```

    * Update the /etc/hosts file

        ```bash
        sudo vim /etc/hosts
        127.0.0.1 jumpbox.localdomain jumpbox localhost localhost.localdomain
        ```

    * Reboot instance

        ```bash
        sudo reboot
        ```

    ***Reference:** <https://aws.amazon.com/premiumsupport/knowledge-center/linux-static-hostname-rhel-centos-amazon/>*

1. Prepare Instances for Replica Set

    * Open /etc/hosts

        ```bash
        sudo vi /etc/hosts
        ```

    * Add IPs of EC2 Instances

        ```bash
        10.0.1.115  primary
        10.0.1.165  secondary1
        10.0.1.175  secondary2
        10.0.1.107  secondary3
        10.0.1.211  secondary4
        ```

    ***Note:** Do it for each instance*

    * Making sure the hostnames are changed

        * Check host name

            ```bash
            sudo hostname -f
            ```

        * Change if not changed yet

            ```bash
            sudo hostnamectl set-hostname <new hostname>
            ```

        * Restart instance after change

            ```bash
            sudo reboot
            ```

1. Initiate Replica-set
    * Open mongo cli in primary

        ```bash
        mongo
        ```

    * Initiate Replica-set

        ```bash
        rs.initiate( {
            _id : "cmpe281",
            members: [
                { _id: 0, host: "primary:27017" },
                { _id: 1, host: "secondary1:27017" },
                { _id: 2, host: "secondary2:27017" },
                { _id: 3, host: "secondary3:27017" },
                { _id: 4, host: "secondary4:27017" }
            ]
        })
        ```

    **Challenge** : Faced some issues connecting to instances in private network.

1. Create Admin Account

    * Open mongo-cli in primary instance

        ```bash
        mongo
        ```

    * Use admin database

        ```bash
        use admin
        ```

    * Create admin account

        ```bash
        db.createUser( {
            user: "admin",
            pwd: "cmpe281",
            roles: [{ role: "root", db: "admin" }]
        });
        ```

1. Open MongoDB cli in instances. From now on, in order to access the mongo-cli use this admin credentials

    ```bash
    mongo -u admin -p cmpe281 --authenticationDatabase admin
    ```

    ***NOTE:** The cluster will choose its new primary when the existing primary instance is down.*

#### Test MongoDB Cluster

1. Add test document into primary

    ```bash
    db.test.save( { a : 1 } )
    ```

1. Find this test document

    ```bash
    db.test.find()
    ```

1. Update test document

    ```bash
    db.test.replaceOne( { a : 1 }, { a : 2 } )
    ```

    All this commands will run properly from the primary node of the cluster.

1. In order to allow queries from secondary, set **Slave OK**

    ```bash
    db.getMongo().setSlaveOk()
    ```

1. Now try finding this document from secondary nodes

    ```bash
    db.test.find()
    ```

---

### **Riak**

#### **Riak Progress**

* [x] Create Riak Cluster
* [x] Test Riak Cluster
* [x] Create Video

***Reference:**<https://github.com/paulnguyen/cmpe281/blob/master/labs/lab4/aws-riak-database-cluster.md>*

#### Create Riak Cluster

1. Launch Instances

    * AMI: Riak KV 2.2 Series
    * Instance Type: t2.micro
    * Number of instances: 5
    * Network: CMPE281
    * Subnet: Private
    * Auto-assign IP: Disable
    * Security Group: Riak
        * Ports: 22, 8087, 8098
    * Keypair: cmpe281-us-west-1.pem

1. Open Additional Security Ports for Cluster

    * Ports :4369, 6000-7999, 8099, 9080
    * Source: sg-0c51891648571d2c7

1. Riak Instances Information

    |Instance|Private IP|SSH|
    |-|-|-|
    |riak1|10.0.1.236|ssh -i "cmpe281-us-west-1.pem" ec2-user@10.0.1.236|
    |riak2|10.0.1.254|ssh -i "cmpe281-us-west-1.pem" ec2-user@10.0.1.254|
    |riak3|10.0.1.233|ssh -i "cmpe281-us-west-1.pem" ec2-user@10.0.1.233|
    |riak4|10.0.1.251|ssh -i "cmpe281-us-west-1.pem" ec2-user@10.0.1.251|
    |riak5|10.0.1.159|ssh -i "cmpe281-us-west-1.pem" ec2-user@10.0.1.159|

    ***Note:** I have given instances names riak1, riak2, riak3, riak4, riak 5 for better understanding.*

1. SSH into Riak instances

    * SSH into **jumpbox**

        ```bash
        chmod 400 cmpe281-us-west-1.pem

        ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-13-56-161-186.us-west-1.compute.amazonaws.com
        ```

    * SSH into **riak1**

        ```bash
        ssh -i "cmpe281-us-west-1.pem" ec2-user@<Private IP of RIAK Instance>
        ```

1. Enable ports for cluster
    * Edit **/etc/riak/riak.conf**

        ```bash
        sudo vi /etc/riak/riak.conf
        ```

1. Connect every instances to **riak1** by executing this command from every other instance

    ```bash
    sudo riak-admin cluster join riak@10.0.1.236
    ```

1. Plan riak cluster

    ```bash
    sudo riak-admin cluster plan
    ```

1. Check the cluster planning status

    ```bash
    sudo riak-admin cluster status
    ```

1. Commit changes of Riak cluster

    ```bash
    sudo riak-admin cluster commit
    ```

1. Check the cluster status

    ```bash
    sudo riak-admin member_status
    ```

#### Testing Riak Cluster

1. Create Bucker

    ```curl
    curl -i http://10.0.1.236:8098/buckets?buckets=true
    ```

1. Write Data

    ```curl
    curl -v -XPUT -d '{"jay":"parekh"}' \
    http://10.0.1.236:8098/buckets/bucket/keys/key1?returnbody=true
    ```

1. Read Data

    ```curl
    curl -i http://10.0.1.254:8098/buckets/bucket/keys/key1
    ```

    Here, we can read data from any node of the cluster since data is replicated in each node.

### Checking Partition Tolerence

To analyse the CAP theorerm for MongoDB and Riak, we have to create a partition in the cluster and the cluster will behave in either CP or AP manner.

* To create a partition, we can stop the communication between two nodes of the cluster.

    ```bash
    sudo iptables -I INPUT -s <IPADDRESS> -j DROP
    ```

* To lift the partition, we can just re-establish communication again.

    ```bash
    sudo iptables -D INPUT -s <IPADDRESS> -j DROP
    ```

#### Checking Partition Tolerence in MongoDB

|Instance|IP|SSH|
|--------|--|---|
|jumpbox|13.56.161.186|ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-13-56-161-186.us-west-1.compute.amazonaws.com|
|mongo-primary|10.0.1.115|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.115|
|mongo-secondary-1|10.0.1.165|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.165|
|mongo-secondary-2|10.0.1.175|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.175|
|mongo-secondary-3|10.0.1.107|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.107|
|mongo-secondary-4|10.0.1.211|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.211|

* MongoDB cli comand

    ```bash
    mongo -u admin -p cmpe281 --authenticationDatabase admin
    ```

* SlaveOk

    ```bash
    db.getMongo().setSlaveOk()
    ```

##### Test 1: Creating partition in slave node

* Stopping communication of **mongo-primary** with others

    ```bash
    sudo iptables -I INPUT -s 10.0.1.165 -j DROP
    sudo iptables -I INPUT -s 10.0.1.175 -j DROP
    sudo iptables -I INPUT -s 10.0.1.107 -j DROP
    sudo iptables -I INPUT -s 10.0.1.11 -j DROP
    ```
* Re-establishing the communication of **mongo-primary** with others

    ```bash
    sudo iptables -D INPUT -s 10.0.1.165 -j DROP
    sudo iptables -D INPUT -s 10.0.1.175 -j DROP
    sudo iptables -D INPUT -s 10.0.1.107 -j DROP
    sudo iptables -D INPUT -s 10.0.1.211 -j DROP
    ```

##### Test 1: Creating partition in master node

* Stopping communication of **mongo-primary** with others

    ```bash
    sudo iptables -I INPUT -s 10.0.1.165 -j DROP
    sudo iptables -I INPUT -s 10.0.1.115 -j DROP
    sudo iptables -I INPUT -s 10.0.1.107 -j DROP
    sudo iptables -I INPUT -s 10.0.1.211 -j DROP
    ```
* Re-establishing the communication of **mongo-primary** with others

    ```bash
    sudo iptables -D INPUT -s 10.0.1.165 -j DROP
    sudo iptables -D INPUT -s 10.0.1.115 -j DROP
    sudo iptables -D INPUT -s 10.0.1.107 -j DROP
    sudo iptables -D INPUT -s 10.0.1.211 -j DROP
    ```

* MongoDB queries

    ```bash
    use project;
    db.getCollectionNames();
    db.createCollection("personal");
    db.personal.insert({test:1});
    db.personal.find();
    ```

#### Checking Partition Tolerence in Riak

|Instance|Private IP|SSH|
|-|-|-|
|riak1|10.0.1.236|ssh -i "cmpe281-us-west-1.pem" ec2-user@10.0.1.236|
|riak2|10.0.1.254|ssh -i "cmpe281-us-west-1.pem" ec2-user@10.0.1.254|
|riak3|10.0.1.233|ssh -i "cmpe281-us-west-1.pem" ec2-user@10.0.1.233|
|riak4|10.0.1.251|ssh -i "cmpe281-us-west-1.pem" ec2-user@10.0.1.251|
|riak5|10.0.1.159|ssh -i "cmpe281-us-west-1.pem" ec2-user@10.0.1.159|

* Stopping communication of **riak1** with others

    ```bash
    sudo iptables -I INPUT -s 10.0.1.254 -j DROP
    sudo iptables -I INPUT -s 10.0.1.233 -j DROP
    sudo iptables -I INPUT -s 10.0.1.251 -j DROP
    sudo iptables -I INPUT -s 10.0.1.159 -j DROP
    ```
* Re-establishing the communication of **riak1** with others

    ```bash
    sudo iptables -D INPUT -s 10.0.1.254 -j DROP
    sudo iptables -D INPUT -s 10.0.1.233 -j DROP
    sudo iptables -D INPUT -s 10.0.1.251 -j DROP
    sudo iptables -D INPUT -s 10.0.1.159 -j DROP
    ```

##### Performing CRUD outside of partition

1. Riak queries

    * Creating bucket

    ```bash
    curl -i http://10.0.1.254:8098/buckets?buckets=true
    ```

    * Inserting key

    ```bash
    curl -v -XPUT -d '{"personal":"project"}' http://10.0.1.254:8098/buckets/bucket/keys/key10?returnbody=true
    ```

    * Fetching key outside partiton

    ```bash
    curl -i http://10.0.1.254:8098/buckets/bucket/keys/key10
    ```

    * Fetching key in partiton
    ```bash
    curl -i http://10.0.1.236:8098/buckets/bucket/keys/key10
    ```

---

## Observations

### Observations in MongoDB

#### In normal condition

MongoDB works on master-slave architecture and so it allows writes to the only master node and not the other nodes. Read permission is given to every slave nodes. When the data is inserted into the master node, it will be replicated to all the other nodes. Each slave node will require a permission to read the data.

If the master node is unreachable, then an election algorithm will elect the new master node based on priority. If the previous node becomes reachable in some way then it will be converted to the secondary node.

### During Partition

If the partition has occurred into the slave node then that particular node will still let users access the database. But the new data inserted into the cluster will not be replicated into the partition node. It will still show the un-updated data.

If the partition has occurred into the master node, then the cluster will run an election algorithm and elect a new master. The old primary will be changed to the secondary node. The new primary will allow the writes and the old primary will only allow reads for the data that it has.

### After Partition Recovery

After the partition has been recovered, the data from the primary will be replicated to the secondary node eventually. And all the data will be consistent.

If the partition has occurred in the primary then after the partition recovery, the primary will be joined into the cluster as a secondary and the data from the new master will be replicated to this node.

### Conclusion for MongoDB

Where the consistency in data is more important than availability, MongoDB is a great option.

### Observations in Riak

#### In normal condition

Riak cluster works as a chain and each cluster node is a part of the chain. Each node allows reads as well as writes. When a write is done in one node, the changes are replicated to all the other nodes.

#### During Partition

Riak allows read and writes during the partition, the partition node will let us change the values.  We can also change the same data in multiple nodes.

#### After Partition

After the partition has been resolved, the latest changes will be replicated to all the nodes across the cluster.

#### Conclusion for Riak

Riak is best used in places where data availability is more important than consistency.

---

### MongoDB Sharding

#### Step 1: Create Security Groups

Name: mongodb-shard-internal
|Port|Source|
|-|-|
|22|Anywhere|
|27017|sg-0291055648f99c5f5|
|27018|sg-0291055648f99c5f5|
|27019|sg-0291055648f99c5f5|

Name: mongodb-shard-internal
|Port|Source|
|-|-|
|22|Anywhere|
|27017|Anywhere|

#### Step 2: Create MongoDB sharded cluster

1. Launch Instance
    ```bash
    AMI: Amazon Linux AMI 2018.03.0 (HVM)
    Type: t2.micro
    VPC: CMPE281
    Subnet: Public Subnet
    Auto-assign Public IP: Enable
    Tags:
        Key: Name
        Tag: mongo-shard-config-1
    Security Group: mongodb-shard-internal
    ```

1. SSH into instance
    ```bash
    chmod 400 cmpe281-us-west-1.pem
    ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-13-56-210-79.us-west-1.compute.amazonaws.com
    ```

1. Install MongoDB

    * Configure the package management system

        * Create mongodb-org-4.0.repo file
        ```bash
        sudo nano /etc/yum.repos.d/mongodb-org-4.0.repo
        ```

        * Insert content into file
        ```bash
        [mongodb-org-4.0]
        name=MongoDB Repository
        baseurl=https://repo.mongodb.org/yum/amazon/2013.03/mongodb-org/4.0/x86_64/
        gpgcheck=1
        enabled=1
        gpgkey=https://www.mongodb.org/static/pgp/server-4.0.asc
        ```

        * Setting MongoDB to start at boot
        ```bash
        sudo chkconfig mongod on
        ```

1. Create Image

    Name: mongo-shard

1. Create another instance
    ```bash
    AMI: mongo-shard
    Type: t2.micro
    VPC: CMPE281
    Subnet: Public Subnet
    Auto-assign Public IP: Enable
    Tags:
        Key: Name
        Tag: mongo-shard-config-1
    Security Group: mongodb-shard-internal
    ```

1. Instances description

    |Instance|IP|SSH|
    |-|-|-|
    |jumpbox|34.219.207.53|ssh -i "cmpe281-us-west-2.pem" ec2-user@ec2-34-219-207-53.us-west-2.compute.amazonaws.com|
    |mongos|10.0.1.239|ssh -i "cmpe281-us-west-2.pem" ec2-user@10.0.1.239|
    |mongo-shard-config-1|10.0.1.78|ssh -i "cmpe281-us-west-2.pem" ec2-user@10.0.1.78|
    |mongo-shard-config-2|10.0.1.65|ssh -i "cmpe281-us-west-2.pem" ec2-user@10.0.1.65|
    |mongo-shard-1.1|10.0.1.220|ssh -i "cmpe281-us-west-2.pem" ec2-user@10.0.1.220|
    |mongo-shard-1.2|10.0.1.44|ssh -i "cmpe281-us-west-2.pem" ec2-user@10.0.1.44|
    |mongo-shard-2.1|10.0.1.141|ssh -i "cmpe281-us-west-2.pem" ec2-user@10.0.1.141|
    |mongo-shard-2.2|10.0.1.129|sssh -i "cmpe281-us-west-2.pem" ec2-user@10.0.1.129|

1. Changes in mongodb-shard-config instances

    *  Make directory for data
        ```bash
        sudo mkdir -p /data/db
        ```

    * Give permission to mongod
        ```bash
        sudo chown -R mongod:mongod /data/db
        ```

    * Change mongod.conf

        * Open file
            ```bash
            sudo vi /etc/mongod.conf
            ```

        * File changes
            ```bash
            storage:
                dbpath: /data/db

            net:
                port: 27019
                bindIp: 0.0.0.0

            replication:
                replSetName: crs

            sharding:
                clusterRole: configsvr
            ```

    * Start mongod daemon with new changes
        ```bash
        sudo mongod --config /etc/mongod.conf --logpath /var/log/mongodb/mongod.log
        ```

1. mongo-shard-conf replicaset.
    ```bash
    mongo -port 27019
    ```
    ```bash
    rs.initiate(
        {
            _id: "crs",
            configsvr: true,
            members: [
                { _id : 0, host : "10.0.1.78:27019"},
                { _id : 1, host : "10.0.1.65:27019"}
            ]
        }
    )

1. Changes in mongodb-shard-1.x instances

    *  Make directory for data
        ```bash
        sudo mkdir -p /data/db
        ```

    * Give permission to mongod
        ```bash
        sudo chown -R mongod:mongod /data/db
        ```

    * Change mongod.conf

        * Open file
            ```bash
            sudo vi /etc/mongod.conf
            ```

        * File changes
            ```bash
            storage:
                dbpath: /data/db

            net:
                port: 27018
                bindIp: 0.0.0.0

            replication:
                replSetName: rs0

            sharding:
                clusterRole: shardsvr
            ```

    * Start mongod daemon with new changes
        ```bash
        sudo mongod --config /etc/mongod.conf --logpath /var/log/mongodb/mongod.log
        ```
1. mongo-shard-1 replicaset.
    ```bash
    mongo -port 27018
    ```
    ```bash
    rs.initiate(
        {
            _id: "rs0",
            members: [
                { _id : 0, host : "10.0.1.220:27018"},
                { _id : 1, host : "10.0.1.44:27018"}
            ]
        }
    )

1. Changes in mongodb-shard-2.x instances

    *  Make directory for data
    ```bash
    sudo mkdir -p /data/db
    ```

    * Give permission to mongod
    ```bash
    sudo chown -R mongod:mongod /data/db
    ```

    * Change mongod.conf

        * Open file
            ```bash
            sudo vi /etc/mongod.conf
            ```

        * File changes
            ```bash
            storage:
                dbpath: /data/db

            net:
                port: 27018
                bindIp: 0.0.0.0

            replication:
                replSetName: rs1

            sharding:
                clusterRole: shardsvr
            ```

    * Start mongod daemon with new changes
        ```bash
        sudo mongod --config /etc/mongod.conf --logpath /var/log/mongodb/mongod.log
        ```

1. mongo-shard-2 replicaset.
    ```bash
    mongo -port 27018
    ```
    ```bash
    rs.initiate(
        {
            _id: "rs1",
            members: [
                { _id : 0, host : "10.0.1.141:27018"},
                { _id : 1, host : "10.0.1.129:27018"}
            ]
        }
    )
    ```

1. Changes in mongose instance
    * Open file
        ```bash
        sudo vi /etc/mongod.conf
        ```

    * File changes
        ```bash
        #storage:

        net:
            port: 27018
            bindIp: 0.0.0.0

        sharding:
            configDB: crs/10.0.1.78, 10.0.1.65:27019
        ```

1. Start mongos
    ```bash
    sudo mongos --config /etc/mongod.conf --logpath /var/log/mongodb/mongod.log
    ```

1. Start mongos-cli
    ```bash
    mongo -port 27017
    ```

1. Use admin database
    ```bash
    use admin
    ```

1. Add shards in mongos
    ```bash
    sh.addShard("rs0/10.0.1.220:27018,10.0.1.44:27018");
    sh.addShard("rs1/10.0.1.141:27018,10.0.1.129:27018");
    ```

### Step 3: Test Sharded Cluster

1. List shards
    ```bash
    db.adminCommand({listShards:1})
    ```

1. Create database
    ```bash
    use test
    ```

1. Enable sharding on database
    ```bash
    db.runCommand({enableSharding: "test"})
    ```
1. Create collection with shard key
    ```bash
    db.runCommand( { shardcollection : "test.users", key : { city : 1} } );
    ```
1. Adding mock data
    ```bash
    db.users.insert("id":1,"first_name":"Esdras","last_name":"Ollander","email":"eollander0@msu.edu","gender":"Male","city":"Miaoya"});
    ```
1. Reading data
    ```bash
    db.users.find();
    ```
