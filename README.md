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
    * Tag: mongo-master
    * Security Group: mongo
      * Ports: 22, 27017
    * Keypair: cmpe281-us-west-1.pem

1. Give Elastic IP to **mongo-master**
    * Allocate New Address
    * Associate it to **mongo-master**
    * Name it as **mongo-master**
    * Elastic IP for **mongo-master**:
        ```bash
        52.9.23.124
        ```

1. Connect to **mongo-master**
     ```bash
    chmod 400 cmpe281-us-west-1.pem
    ssh -i "cmpe281-us-west-1.pem" ec2-user@ec2-52-9-23-124.us-west-1.compute.amazonaws.com
    ```

1. Install MongoDB Community Edition
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
    * Install MongoDB packages.
        ```bash
        sudo yum install -y mongodb-org
        ```