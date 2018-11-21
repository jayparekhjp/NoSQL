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

---

## **Journal**

---

### **MongoDB**

#### **Mongo Progress**

* [x] Create MongoDB Cluster
* [x] Test MongoDB Cluster
* [ ] Test CP Properties
* [ ] Create MongoDB Shards
* [ ] Create GO API

#### Create MongoDB Cluster

***Reference**: <https://github.com/paulnguyen/cmpe281/blob/master/labs/lab4/aws-mongodb-replica-set.md>*

1. Create Jumpbox

    **Note:** *Since all the instances will be in private subnet, jumpbox is needed to access them*
    * AMI: Amazon Linux AMI 2018.03.0 (HVM), SSD Volume Type
    * Instance Type: t2.micro
    * Network: CMPE281
    * Subnet: Public Subnet
    * Auto-assign Public IP: Enable
    * Tag: jumpbox
    * Security Group: **jumpbox**
      * Ports: 22, 80
    * Keypair: cmpe281-us-west-1.pem

1. Creating an EC2 Instance
    * AMI: Ubuntu Server 16.04 LTS (HVM), SSD Volume Type
    * Instance Type: t2.micro
    * Network: CMPE281
    * Subnet: Private Subnet
    * Auto-assign Public IP: Disable
    * Tag: **mongo-primary**
    * Security Group: mongo
      * Ports: 22, 27017
    * Keypair: cmpe281-us-west-1.pem

1. Connecting to **mongo-primary**

    * Upload key to **jumpbox**
        ```bash
        scp -i cmpe281-us-west-1.pem cmpe281-us-west-1.pem ec2-user@ec2-13-56-16-49.us-west-1.compute.amazonaws.com:
        ```
    * Connect to **jumpbox**
        ```bash
        chmod 400 cmpe281-us-west-1.pem
        ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-user@ec2-13-56-16-49.us-west-1.compute.amazonaws.com
        ```
    * Connect to **mongo-primary**
        ```bash
        chmod 400 cmpe281-us-west-1.pem
        ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.67
        ```

1. Install MongoDB

    ***Note:*** *Start the NAT-gateway instance of the VPC with Elastic IP in order to provide internet access to private subnet instances.*

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
    |mongo-primary|10.0.1.115|ssh -i "cmpe281-us-west-1.pem" root@ec2-54-183-146-72.us-west-1.compute.amazonaws.com|
    |mongo-secondary-1|10.0.1.165|ssh -i "cmpe281-us-west-1.pem" root@ec2-13-56-59-10.us-west-1.compute.amazonaws.com|
    |mongo-secondary-2|10.0.1.175|ssh -i "cmpe281-us-west-1.pem" root@ec2-18-144-45-78.us-west-1.compute.amazonaws.com|
    |mongo-secondary-3|10.0.1.107|ssh -i "cmpe281-us-west-1.pem" root@ec2-18-144-34-186.us-west-1.compute.amazonaws.com|
    |mongo-secondary-4|10.0.1.211|ssh -i "cmpe281-us-west-1.pem" root@ec2-54-219-185-196.us-west-1.compute.amazonaws.com|

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
                { _id: 4, host: "secondary4:27017" },
                { _id: 5, host: "secondary5:27017" }
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
    * From now on, in order to access the mongo-cli use this admin credentials
        ```bash
        mongo -u admin -p cmpe281 --authenticationDatabase admin
        ```

    ***NOTE:** The cluster will choose its new primary when the existing primary instance is down.*

1. Test cluster by adding test data into master.

    * Add test document into primary
        ```bash
        db.test.save( { a : 1 } )
        ```
    * Find this test document
        ```bash
        db.test.find()
        ```
    * Update test document
        ```bash
        db.test.replaceOne( { a : 1 }, { a : 2 } )
        ```
    All this commands will run properly from the primary node of the cluster.

    * In order to allow queries from secondary, set **Slave OK**
        ```bash
        db.getMongo().setSlaveOk()
        ```
    * Now try finding this document from secondary nodes
        ```bash
        db.test.find()
        ```

---

### **Cassandra**

* [x] Create Kubernetes Cluster on AWS
* [ ] Install Cassandra on Kubernetes Cluster
* [ ] Test CP Properties

<!--
    #### **Cassandra Progress**

    * [ ] Create Cassandra Cluster
    * [ ] Test CP Properties
    * [ ] Test Cassandra Extras
    * [ ] Create Shards

    #### Create Cassandra Cluster

    1. Launch EC2 Instance

        * AMI: Amazon Linux AMI 2018.03.0 (HVM), SSD Volume Type
        * Instance Type: t2.micro
        * Network: CMPE281
        * Subnet: Private Subnet
        * Auto-assign Public IP: Disable
        * Tag: **cassandra-1**
        * Security Group: cassandra
        * Ports: 22, 80, 9042
        * Keypair: cmpe281-us-west-1.pem

    1. Connecting to **cassandra-1**

        * Connect to **jumpbox**
            ```bash
            ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-13-56-16-49.us-west-1.compute.amazonaws.com
            ```

        * Connect to **cassandra-1**
            ```bash
            ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.82
            ```

    1. Install Java JVM

        * Add Personal Package Archives
            ```bash
            sudo add-apt-repository ppa:webupd8team/java
            ```
        * Update the package database
            ```bash
            sudo apt-get update
            ```
        * Install the Oracle JRE
            ```bash
            sudo apt-get install oracle-java8-set-default
            ```
        * Verify
            ```bash
            java -version
            ```

    1. Install Cassandra

        * Add the Cassandra Repository
            ```bash
            echo "deb http://www.apache.org/dist/cassandra/debian 311x main" | sudo tee -a /etc/apt/sources.list.d/cassandra.sources.list
            ```
        * Add the Cassandra Repository Keys
            ```bash
            curl https://www.apache.org/dist/cassandra/KEYS | sudo apt-key add -
            ```
        * Update the package index
            ```bash
            sudo apt-get update
            ```
        * Install Apache Cassandra
            ```bash
            sudo apt-get install cassandra
            ```
        * Check the Status of the Apache Cassandra Service
            ```bash
            sudo systemctl status cassandra.service
            ```
        * Start the Apache Cassandra Service
            ```bash
            sudo systemctl start cassandra.service
            ```
        * Stop the Apache Cassandra Service
            ```bash
            sudo systemctl stop cassandra.service
            ```
        * Enable Apache Cassandra Service on System Boot
            ```bash
            sudo systemctl enable cassandra.service
            ```
        ***Reference:** <https://www.rosehosting.com/blog/how-to-install-apache-cassandra-on-ubuntu-16-04/>*
-->

#### Kubernetes Setup Local

***Note:** Trying to deploy cassandra on AWS EKS.*

1. Install **kops**

    *Used to setup infrastructure in AWS*
    ```bash
    curl -LO https://github.com/kubernetes/kops/releases/download/1.8.1/kops-linux-amd64 

    sudo mv kops-linux-amd64 /usr/local/bin/kops && sudo chmod a+x /usr/local/bin/kops
    ```

1. Install **kubectl**

    *To interact with Kubernetes cluster in AWS*
    ```bash
    curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.7.16/bin/linux/amd64/kubectl

    sudo mv kubectl /usr/local/bin/kubectl && sudo chmod a+x /usr/local/bin/kubectl
    ```

1. Install **awscli**

    ***kops** uses awscli to interact with AWS*
    ```bash
    sudo apt-get install awscli
    ```

1. Create an IAM user
    ```text
    User name: jpCassandra
    Access type: Programmatic Access, AWS Management Console access
    Attach existing policies directly: AmazonEC2FullAccess, AmazonRoute53FullAccess, AmazonS3FullAccess, IAMFullAccess, AmazonVPCFullAccess
    ```

1. Configure AWS from terminal
    ```bash
    aws configure
    ```

1. Create S3 Bucket
    ```bash
    aws s3api create-bucket --bucket jp-kops-cassandra --region eu-west-1 --create-bucket-configuration LocationConstraint=eu-west-1
    ```

1. Generate a Public/Private key-pair

    *This key-pair will be used to access the EC2 instances*
    ```bash
    ssh-keygen -f jp-kops-cassandra
    ```

1. Cluster Definition
    ```bash
    kops create cluster \
    --cloud=aws \
    --name=jp-kops-cassandra.k8s.local \
    --zones=eu-west-1a,eu-west-1b,eu-west-1c \
    --master-size="t2.micro" \
    --master-zones=eu-west-1a,eu-west-1b,eu-west-1c \
    --node-size="t2.micro" \
    --ssh-public-key="jp-kops-cassandra.pub" \
    --state=s3://jp-kops-cassandra \
    --node-count=6
    ```

1. Apply the cluster definition
    ```bash
    kops update cluster --name=jp-kops-cassandra.k8s.local --state=s3://jp-kops-cassandra --yes
    ```

1. Check the Kubernetes master nodes
    ```bash
    kubectl get no -L failure-domain.beta.kubernetes.io/zone -l kubernetes.io/role=master
    ```

1. Check the Kubernetes nodes
    ```bash
    kubectl get no -L failure-domain.beta.kubernetes.io/zone -l kubernetes.io/role=node
    ```

1. Destroy the environment
    ```bash
    kops delete cluster --name=jp-kops-cassandra.k8s.local --state=s3://jp-kops-cassandra --yes
    ```