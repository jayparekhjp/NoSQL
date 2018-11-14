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

## Journal

### **MongoDB**

#### Step 1: Create MongoDB Cluster

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
        scp -i cmpe281-us-west-1.pem cmpe281-us-west-1.pem ec2-user@ec2-18-144-42-185.us-west-1.compute.amazonaws.com:
        ```
    * Connect to **jumpbox**
        ```bash
        chmod 400 cmpe281-us-west-1.pem
        ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-18-144-42-185.us-west-1.compute.amazonaws.com
        ```
    * Connect to **mongo-primary**
        ```bash
        chmod 400 cmpe281-us-west-1.pem
        ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.115
        ```

1. Install MongoDB

    **Note:** *Start the NAT-gateway instance of the VPC with Elastic IP in order to provide internet access to private subnet instances.*
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
    <!--1. Launch MongoDB as a service
        ```bash
        sudo vim /etc/systemd/system/mongodb.service
        ```
    -->
    <!--1. Install MongoDB for Amazon Linux
        * Configure the package management system.
            * Create a /etc/yum.repos.d/mongodb-org-4.0.repo file to install MongoDB directly using yum.
                ```bash
                sudo vi /etc/yum.repos.d/mongodb-org-4.0.repo
                ```
            * File Content:
                ```bash
                [mongodb-org-4.0]
                name=MongoDB Repository
                baseurl=https://repo.mongodb.org/yum/amazon/2013.03/mongodb-org/4.0/x86_64/
                gpgcheck=1
                enabled=1
                gpgkey=https://www.mongodb.org/static/pgp/server-4.0.asc
                ```
        * Install MongoDB packages
            ```bash
            sudo yum install -y mongodb-org
            ```
    -->
    <!--1. Run Mongo Commands to Test Installation
        * Verify mongod process has started
            ```bash
            sudo cat /var/log/mongodb/mongod.log
            ```
        * Ensure MongoDB will start after reboot also
            ```bash
            sudo chkconfig mongod on
            ```
        * Check Status MongoDB
            ```bash
            sudo service mongod status
            ```
        * Start MongoDB
            ```bash
            sudo service mongod start
            ```
        * Stop MongoDB
            ```bash
            sudo service mongod stop
            ```
        * Restart MongoDB
            ```bash
            sudo service mongod restart
            ```
        * Begin MongoDB CLI
            ```bash
            mongo
            ```
        * Exit MongoDB CLI
            ```bash
            exit
            ```
    -->

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
    |mongo-primary|10.0.1.115|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.115|
    |mongo-secondary-1|10.0.1.95|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.95|
    |mongo-secondary-2|10.0.1.68|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.68|
    |mongo-secondary-3|10.0.1.212|sssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.212|
    |mongo-secondary-4|10.0.1.199|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.199|
    |mongo-secondary-5|10.0.1.30|ssh -i "cmpe281-us-west-1.pem" ubuntu@10.0.1.30|

