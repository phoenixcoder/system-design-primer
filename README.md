# Introduction to Large-Scale System Design

## Purpose
This is a hands-on introduction to the basic components involved in the design of large-scale systems.
The examples here are backed by AWS Cloudformation Templates, which automates their setup to let you
play/experiment with the components. I'll be adding more examples here over time, on different platforms e.g. GCP, Azure, etc. to address certain problem scenarios, and ways we can solve them. I have and will continue to design all these examples with the free-tier and cost in mind of any platform we use. 

## Before You Get Started
0. Know the basics of Git: TODO Link git tutorial here.
1. Create an AWS account: [AWS Documentation: How do I create and activate a new Amazon Web Services account?](https://aws.amazon.com/premiumsupport/knowledge-center/create-and-activate-aws-account/)
2. [Create an AWS Admin User, and Locally Save Access Key ID and Secret Access Key.](TODO Update with Section Anchor Link)
3. [Configure AWS CLI with the Administrative User Credentials.](TODO Update with Section Anchor Link) 
4. CAUTION, I attempt to contain everything within a Cloudformation template that usually keeps track of what resources it
creates and what resources it deletes when removed. However, if you run into any problems with deletion, just know
if some resources stick around e.g. EC2 instances, that once the free tier expires you may be charged.
5. All of the following was tested using the following:
 - macOS Mojave 10.14.5
 - iTerm2 Build 3.2.9

 In general, it's probably sufficient to just have some version of macOS that's 10.x.x+ and vanilla terminal.
 For Windows users, your time is coming for vetted instructions.
6. For Mac users, you will need the following:
 - [AWS Documentation: Install the AWS CLI on macOS Using pip](https://docs.aws.amazon.com/cli/latest/userguide/install-macos.html#awscli-install-osx-pip)
 - [AWS Documentation: Add the AWS CLI Executable to Your macOS Command Line Path](https://docs.aws.amazon.com/cli/latest/userguide/install-macos.html#awscli-install-osx-path)

## Basic Distributed Web System Example
### Why is this important?
A distributed web system shows its value when it handles a number of requests beyond what a single computer
can handle. These requests can be anything from just visiting a webpage to downloading a file to borrowing
some time on your system to calculate the largest prime - you never know. If you have these happening thousands
to tens of thousands or more per second, a single computer will probably start to have a lot of trouble completing
all of them successfully. If that single computer can't handle the load, why not give it some help? This example
shows how the basic components work together to give that computer some companions to bear the load.

### What does it look like?
TODO Insert Diagram Here

### Components
- Load Balancer
- Computer Pool
- Web Server Application

### Ready to Play?
#### Running the Example
1. Open your favorite terminal.
1. Clone this repo by running: `git clone https://github.com/phoenixcoder/IntroSystemDesign.git`
1. Run `cd IntroSystemDesign`
1. [Verify Your AWS CLI is Configured Properly](TODO Insert Anchor Section Link)
1. Open the [AWS Cloudformation console](https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks).
1. Keep the [AWS Cloudformation console](https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks) open, and switch back to your terminal to run:
`aws cloudformation create-stack --stack-name basic-distributed-web-system-stack --template-body file://aws/cloudformation/templates/basic-distributed-web-system.json`

   You should get output similar to:
   ```
   {
    "StackId": "arn:aws:cloudformation:us-west-2:<Your Account Number>:stack/basic-distributed-web-system-stack/12345678-90ab-cdef-ghij-klmnopqrstuv"
   }
   ```

   When you go back to the [AWS Cloudformation console](https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks), you should see
   under **Stacks** an item labeled **basic-distributed-web-system-stack**.
1. Wait until the item **basic-distributed-web-system-stack** on the [AWS Cloudformation console](https://us-west-2.console.aws.amazon.com/cloudformation/home?region=us-west-2#/stacks) reports under it **CREATE_COMPLETE**. If it reports **ROLLBACK_IN_PROGRESS** or **ROLLBACK_COMPLETE**, then something has gone wrong, please submit an issue [here](TODO Insert Github Bug Page).

#### Experimentation Suggestions
##### Smash Curl the Load-Balancer
1. Open your favorite terminal.
1. Prepare yourself for the copy-n-pasting of some Bash voodoo.
1. Run the following:
```
export LbDNS=$(aws elb describe-load-balancers --load-balancer-names "Basic-Web-System-LB" | grep DNSName | tr -d " ,\"" | cut -d':' -f2)
```
1. Run the following (a couple of times, if it pleases you):
```
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
2. On the homepage, use the **Find Services** textbox, type in "IAM". Click on the **IAM** option in the drop-down.
3. On the left-panel, click on **Users**.
4. On the right-panel, at the top, click on the **Add User** button.
5. Under **Set user details**, in the **User name** field, type in "admin".
6. Under **Select AWS access type**, next to **Access type**, check the **Programmatic access** checkbox.
7. On the bottom-right, click the **Next: Permissions** button.
8. Under **Set permissions**, click on **Attach existing policies directly**.
9. In the **Search** box, type in "AdministratorAccess" and check the **AdministratorAccess** checkbox.
10. On the bottom-right, click on **Next: Tags**.
11. On the bottom-right, click on **Next: Review**.
12. On the bottom-right, click on **Create: user**.
13. On the very right, under **Secret access key**, click on the **Show** link.
14. Copy the **Access key ID** and **Secret access key**. Put it somewhere safe. If you close the window or skip this step for whatever reason, you will have to create a new access key and secret.
15. At the bottom, click **Close**.

## Configure AWS CLI with Admin User's Access Key ID and Secret Key
1. Open your favorite terminal.
2. Run `aws configure`. The following sequence of lines should appear.
   ```
   AWS Access Key ID [None]: <Your Access Key Here>
   AWS Secret Access Key [None]: <Your Access Secret Here>
   Default region name [None]: us-west-2
   Default output format [None]: json
   ```
3. [Verify Your AWS CLI is Configured Properly](TODO Insert Anchor Section Link)

## Verify Your AWS CLI is Configured Properly
Run `aws iam get-user`, you should get output similar to the following:

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
- [Learn Git Branching](https://learngitbranching.js.org)
- [Git It Electron](https://github.com/jlord/git-it-electron#what-to-install)
