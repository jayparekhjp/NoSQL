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

1. Creating an EC2 Instance
    * AMI: Amazon Linux AMI 2018.03.0 (HVM), SSD Volume Type
    * Instance Type: t2.micro
    * Network: CMPE281
    * Subnet: Public Subnet
    * Auto-assign Public IP: Disable
    * Tag: mongo-primary
    * Security Group: mongo
      * Ports: 22, 27017
    * Keypair: cmpe281-us-west-1.pem

1. Give Elastic IP to **mongo-primary**
    * Allocate New Address
    * Associate it to **mongo-primary**
    * Name it as **mongo-primary**
        ```bash
        52.9.23.124
        ```

1. Connect to **mongo-primary**
     ```bash
    chmod 400 cmpe281-us-west-1.pem
    ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-52-9-23-124.us-west-1.compute.amazonaws.com
    ```

1. Install MongoDB
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
1. Run Mongo Commands to Test Installation
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
1. Create Image of **mongo-primary**
    * Image Name: mongodb
    * Image Description: mongodb-version 4

1. Launch Secondary Instances
    * AMI: mongodb
    * Instance Type: t2.micro
    * Number of Instances: 5
    * Network: CMPE281
    * Subnet: Public Subnet
    * Auto-assign Public IP: Enable
    * Security Group: mongo
        * Ports: 22,27017
    * Key: cmpe281-us-west-2.pem
    * *Give them names mongo-secondary1, mongo-secondary2, mongo-secondary3, mongo-secondary4, mongo-secondary5 for better understanding*

1. Information of Instances

    |Instance|IP|SSH|
    |--------|--|---|
    |mongo-primary|52.9.23.124|ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-52-9-23-124.us-west-1.compute.amazonaws.com
    |mongo-secondary1|52.53.174.32|ssh -i "cmpe281-us-west-1.pem" root@ec2-52-53-174-32.us-west-1.compute.amazonaws.com
    |mongo-secondary2|13.56.178.86|ssh -i "cmpe281-us-west-1.pem" root@ec2-13-56-178-86.us-west-1.compute.amazonaws.com
    |mongo-secondary3|18.144.11.254|ssh -i "cmpe281-us-west-1.pem" root@ec2-18-144-11-254.us-west-1.compute.amazonaws.com
    |mongo-secondary4|13.57.178.76|ssh -i "cmpe281-us-west-1.pem" root@ec2-13-57-178-76.us-west-1.compute.amazonaws.com
    |mongo-secondary5|13.57.185.229|ssh -i "cmpe281-us-west-1.pem" root@ec2-13-57-185-229.us-west-1.compute.amazonaws.com