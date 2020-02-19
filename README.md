[![Actions Status](https://github.com/piyushsonigra/aws_ipadd/workflows/Build%20&%20Release/badge.svg)](https://github.com/piyushsonigra/aws_ipadd/actions)


# aws_ipadd


## OS Support

Currently aws_ipadd supports the following Operating System

- Mac OS X (64bit)
- Linux (64bit)

## :rocket: Installation

Download aws_ipadd for your operating system

  Linux

  ```console
  $ wget -c https://github.com/piyushsonigra/aws_ipadd/releases/latest/download/aws_ipadd_linux_x64.tar.gz -O - | tar -xz -C /usr/local/bin/
  ```

  OSX

  ```console
  $ wget -c https://github.com/piyushsonigra/aws_ipadd/releases/latest/download/aws_ipadd_osx_x64.tar.gz -O - | \ tar -xz -C /usr/local/bin/
  ```

Note: If you get errors related to permission or access, Please run command with `sudo`.

## configuration

Run below commands to conifgure aws_ipadd command. 

  Create directory `~/.aws_ipadd` at your home directory.

  ```console
  $ mkdir ~/.aws_ipadd
  ```
  
  Create configuration file `aws_ipadd` inside `~/.aws_ipadd`.
  
  ```console
  $ touch ~/.aws_ipadd/aws_ipadd
  ```

  Edit the `~/.aws_ipadd/aws_ipadd` file and add below Informations as shown in sample configuration file.
  
  - aws_ipadd profile name in []: 
  `my_project_mysql` and `my_project_ssh` is aws_ipadd profiles to identify configuration which security group rule need to update with port, IP, rule_name and security group region for different AWS account profiles.

  - aws_profile: 
    aws_profile is name of AWS profile configured for awscli.

  - region_name: 
    AWS region name in which security group is present.

  - security_group_id: 
    AWS security group id.

  - rule_name: 
    AWS security group rule name to identify rule purpose.

  - port: 
    Network port to whitelist with IP.
   
  Below is the sample configuration of `~/.aws_ipadd/aws_ipadd` file.
   
  ```console
  $ cat ~/.aws_ipadd/aws_ipadd
  [my_project_ssh]
  aws_profile = my_project
  security_group_id = sg-d26fdre9d
  port = 22
  rule_name = my_office_ssh
  region_name = us-east-1
  
  [my_project_mysql]
  aws_profile = my_project
  security_group_id = sg-dfg9dwe
  port = 3306
  rule_name = my_office_mysql
  region_name = us-east-1
  ```
  
## Usage

Run the aws_ipadd command with aws_ipadd profile.

  ```console
  $ aws_ipadd my_project_ssh
    Your IP 12.10.1.14/32 and Port 22 is whitelisted successfully.
  ```
  
  If your public IP is changed, aws_ipadd will update aws security group rule with your current public IP.
  
  ```console
  $ aaws_ipadd my_project_ssh
    Modifying existing rule...
    Removing old whitelisted IP '12.10.1.14/32'.
    Whitelisting new IP '131.4.10.16/32'.
    Rule successfully updated!
  ```
  
## Licence

- [aws_ipadd](https://github.com/piyushsonigra/aws_ipadd/blob/master/LICENSE)

  
   


