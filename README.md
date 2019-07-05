# Introduction to Large-Scale System Design

## Purpose

This is a hands-on introduction to the components involved in large-scale systems. This guide hopes to help
you learn about the design of such systems by helping you setup and play with the components directly. I'll tackle different problem scenarios and try to build upon the lessons of the previous examples as we go along. The intention is for the scenarios to get more complicated, and to think of algorithms in a distributed manner. If you have any suggestions as to scenarios you want to see addressed, please let me know through the [repo's issues page](https://github.com/phoenixcoder/IntroSystemDesign/issues) or send me a pull request with changes you've made.

These examples are backed by [AWS Cloudformation Templates](#what-is-an-aws-cloudformation-template) that automate infrastructure setup so you can get to playing quicker. I plan to add more examples over time that cross different platforms like GCP and Azure. I have and will continue to design all examples to work within the free-tier of any platform we use.

## Before You Get Started

1. Know Git Basics. Here are some [tutorials](#git-basics-and-tutorials) to get you started.
1. Create an AWS account: [AWS Documentation: How do I create and activate a new Amazon Web Services account?](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/)
1. [Create an AWS Admin User, and Locally Save Access Key ID and Secret Access Key.](#create-an-aws-admin-user-and-locally-save-access-key-id-and-secret-access-key)
1. [Install AWS CLI](#install-aws-cli)
1. [Configure AWS CLI with the Administrative User Credentials](#configure-aws-cli-with-admin-users-access-key-id-and-secret-key)

Tested using the following:

- macOS Mojave 10.14.5
- iTerm2 Build 3.2.9

In general, it's probably sufficient to just have some version of macOS that's 10.x.x+ and vanilla terminal.

For Windows users, your time is coming for vetted instructions.

## Basic Distributed Web System Scenario

[(Nah, don't lecture me, just let me play!)](#ready-to-play)

### Why is this important?

A distributed web system is valuable for handling a large request volume beyond what a single computer
can handle. These requests can be a page visit, to borrowing some time for finding aliens - you never know.
If you have these happening thousands to tens of thousands or more per second, a single computer is unlikely to handle all of it. So, the goal for this design is to find some way to give this computer some friends.

### What does it look like?

![](https://github.com/phoenixcoder/IntroSystemDesign/blob/master/resources/img/basic-distributed-web-system.png)

#### Components

- **Load Balancer:** This component sits in front of a pool of servers, and decides who gets what. Specifically, it takes care of:

  - Distributing network requests across multiple servers.
  - Keeping track of your healthy servers aka servers that can respond to a request, and sending requests only to the healthy ones.
  - Making it really dang easy to add or subtract servers from a pool of computers.

  In this scenario, these load balancers are servers themselves performing and maintaining the properties above.

- **Computer Pool:** These are just your typical plain ol' servers. These have been configured to be on the same subnet, and configurations made so they can see each other. It's useful if you want the servers to talk to one another, which we'll get into in later examples.
- **Web Server Application:** This is like any other web application that services HTTP requests. In this scenario, I've developed a [simple web server](https://github.com/phoenixcoder/IntroSystemDesign/blob/master/go/src/simple_web_server.go) using golang that will respond to requests with a simple greeting as to which server responded.

#### Relationships
Each instance that enters into the **Computer Pool** is required to register with the **Load Balancer** before it can begin operations. Whenever the **Load Balancer** receives a request, it uses a routing algorithm to send a request to one of the servers in the **Computer Pool**, which can by any of the following:

- **Round Robin (yummmm!):** The load is distributed in sequential order. If you had 3 servers, the next four requests would run around the triangle: *1 -> 2 -> 3 -> 1*.
- **Least Connections:** The load balancer keeps a sorted list by least number of connections. It distributes the next request to the one doing the least work.
- **Hashing:** There's some fancy hash algorithm that maps the next request to its intended destination. We can tackle the algorithms for this at a later date.

In this scenario, the **Load Balancer** is an intermediary between the client who sent the request and the server that's serving it. So, remember that there are two sides of the connection to consider here:
- Client <-> Load Balancer
- Load Balancer <-> Web Server Application

### Ready to Play?

#### Running the Example

1. Open your favorite terminal.
1. Clone this repo by running: 
   ```
   git clone https://github.com/phoenixcoder/IntroSystemDesign.git
   ```
1. Run `cd IntroSystemDesign`
1. Make sure you've [configured your AWS CLI](#configure-aws-cli-with-admin-users-access-key-id-and-secret-key).
1. Open the [AWS Cloudformation console](https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks).
1. Keep the [AWS Cloudformation console](https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks) open, and switch back to your terminal to run:
   ```
   aws cloudformation create-stack --stack-name basic-distributed-web-system-stack \
                                   --template-body file://aws/cloudformation/templates/basic-distributed-web-system.json
   ```

   You should get output similar to:

   ```
   aws cloudformation create-stack --stack-name basic-distributed-web-system-stack \
                                   --template-body file://aws/cloudformation/templates/basic-distributed-web-system.json
   {
    "StackId": "arn:aws:cloudformation:us-west-2:<Your Account Number>:stack/basic-distributed-web-system-stack/12345678-90ab-cdef-ghij-klmnopqrstuv"
   }
   ```

   Go back to the [AWS Cloudformation console](https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks), you should see
   under **Stacks** an item labeled **basic-distributed-web-system-stack**.

1. Wait until the item **basic-distributed-web-system-stack** on the [AWS Cloudformation console](https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks) reports under it **CREATE_COMPLETE**. If it reports **ROLLBACK_IN_PROGRESS** or **ROLLBACK_COMPLETE**, then something has gone wrong, please submit an issue [here](https://github.com/phoenixcoder/IntroSystemDesign/issues).

#### Suggested Experiments

##### Smash Curl the Load-Balancer

1. Open your favorite terminal.
1. Prepare yourself for the copy-n-pasting of some Bash voodoo.
1. Run the following (a couple of times, if it pleases you):

   ```
   export LbDNS=$(aws elb describe-load-balancers --load-balancer-names "Basic-Web-System-LB" | grep DNSName | tr -d " ,\"" | cut -d':' -f2)
   for i in {1..10}; do curl -s "http://$LbDNS/request" | grep i-; done | sort | uniq -c
   ```

   You should get results that look like the following:

   ```
   for i in {1..10}; do curl -s "http://$LbDNS/request" | grep i-; done | sort | uniq -c
      5 My name is, i-0b63d29a514360e19.
      5 My name is, i-0fea5dcd64c116cb0.
   ```

   The above means that at least two different servers responded to all your requests.

#### Clean-up

1. Open the [AWS Cloudformation console.](https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks)
1. Click on the item **basic-distributed-web-system-stack**, and in the upper-right corner click **Delete**.
1. Under the item name, the status should change from **CREATE_COMPLETED** to **DELETION_IN_PROGRESS**.
1. Wait for status to change from **DELETION_IN_PROGRESS** to **DELETION_COMPLETED**. Then refresh the page to make sure item has disappeared from **Stacks** panel on left.

# Appendix

## Create an AWS Admin User and Locally Save Access Key ID and Secret Access Key

1. Create and/or login to AWS console in the us-west-2 region: [us-west-2.console.aws.amazon.com](https://us-west-2.console.aws.amazon.com)
1. On the homepage, use the **Find Services** textbox, type in "IAM". Click on the **IAM** option in the drop-down.
1. On the left-panel, click on **Users**.
1. On the right-panel, at the top, click on the **Add User** button.
1. Under **Set user details**, in the **User name** field, type in "admin".
1. Under **Select AWS access type**, next to **Access type**, check the **Programmatic access** checkbox.
1. On the bottom-right, click the **Next: Permissions** button.
1. Under **Set permissions**, click on **Attach existing policies directly**.
1. In the **Search** box, type in "AdministratorAccess" and check the **AdministratorAccess** checkbox.
1. On the bottom-right, click on **Next: Tags**.
1. On the bottom-right, click on **Next: Review**.
1. On the bottom-right, click on **Create: user**.
1. On the very right, under **Secret access key**, click on the **Show** link.
1. Copy the **Access key ID** and **Secret access key**. Put it somewhere safe. If you close the window or skip this step for whatever reason, you will have to create a new access key and secret.
1. At the bottom, click **Close**.

## Install AWS CLI

For Mac users, you will need the following:

- [AWS Documentation: Install the AWS CLI on macOS Using pip](https://docs.aws.amazon.com/cli/latest/userguide/install-macos.html#awscli-install-osx-pip)
- [AWS Documentation: Add the AWS CLI Executable to Your macOS Command Line Path](https://docs.aws.amazon.com/cli/latest/userguide/install-macos.html#awscli-install-osx-path)

## Configure AWS CLI with Admin User's Access Key ID and Secret Key

1. Open your favorite terminal.
1. Run `aws configure`. The following sequence of lines should appear.
   ```
   AWS Access Key ID [None]: <Your Access Key Here>
   AWS Secret Access Key [None]: <Your Access Secret Here>
   Default region name [None]: us-west-2
   Default output format [None]: json
   ```
1. Run `aws iam get-user`, you should get output similar to the following:
   ```
   {
     "User": {
    	    "Path": "/",
          "UserName": "admin",
    	    "UserId": "<Your Access Key>",
          "Arn": "arn:aws:iam::<Your Account Number>:user/admin",
          "CreateDate": "2019-01-01T00:00:00Z"
 	   }
   }
   ```

This validates your configuration is working.

## What is an AWS Cloudformation Template?

A cloudformation template is a document that describes your very own private cloud. The template
and the system that runs it is know as Infrastructure-as-a-Server (IaaS). It automates setting
up small or large infrastructural projects in a shorter period of time vs say buying the hardware
yourself and setting up your own datacenter. If you know what you're looking to do, IaaS can be
a very powerful tool - letting you setup global web services with a team as small as 2 people.

## Git Basics and Tutorials

- [Atlassian Git Tutorials](https://www.atlassian.com/git)
- [Learn Git Branching](https://learngitbranching.js.org)
- [Git It Electron](https://github.com/jlord/git-it-electron#what-to-install)
