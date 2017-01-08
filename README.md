# Maya Server

> Maya exposes its APIs here

A service exposing openebs APIs. 

## Use Cases

### serving EBS needs

Maya server can adapt itself as an Amazon EBS server. This makes it super simple for 
existing EBS clients to talk to Maya `acting as an EBS Server`. One can use existing
docker ebs volume driver or k8s ebs volume persistence plugins talk to maya with zero
learning curve. In other words, infrastructure admins or developers need not code 
anything to interact with `openebs`.

Among other things, this provides an ability to switch from Amazon EBS to OpenEBS & 
vice-versa. This provides the flexibility to use above mix in dev and/or production
environments.
