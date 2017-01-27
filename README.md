# Maya API Server (Work In Progress)

> Maya exposes its APIs here

A service exposing `Elastic Block Store` i.e. EBS APIs, thus making openebs 
storage compatible with EBS APIs.

## Use Cases

### Serving EBS compatibility

Maya server can adapt itself as an Amazon EBS server. This makes it super simple
for existing EBS clients to talk to Maya `i.e. latter acting as an EBS Server`. 
One can use existing docker ebs volume driver or k8s ebs volume persistence 
plugins talk to maya with zero learning curve. In other words, infrastructure 
admins or developers need not code anything to interact with `openebs`.

Among other things, this provides an ability to switch from Amazon EBS to OpenEBS 
& vice-versa. This provides the flexibility to use above mix in dev and/or 
production environments. In other words path towards devops adoption.

## Note

> This is very much a work in progress. Once the code base executes few of the 
mentioned features, the WIP tag will be removed.

## Development of Mayaserver

> These are some of the steps to start off with development.

- git clone https://github.com/openebs/mayaserver.git
- cd to above cloned folder i.e mayaserver
  - vagrant up
  - vagrant ssh
- from within the vagrant VM run below steps:
  - make init
  - make
  - make release
  - sudo nohup mayaserver up &>mserver.log &

## Setting up EBS clients

### Clients that use aws-sdk-go

> Once the credentials are configured, `aws-sdk-go` can find them and use them 
automatically. You don’t need to explicitly reference the credentials in your
code at all. 

- Read through [configuring credentials section](https://github.com/aws/aws-sdk-go#configuring-credentials)
