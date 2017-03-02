# Maya API Server (Work In Progress)

> Maya exposes its APIs here

A service exposing `Elastic Block Store` i.e. EBS APIs, thus making openebs 
storage compatible with EBS APIs.

## Use Cases

### Serving EBS compatibility

Maya server adapts itself like an Amazon EBS server. This makes it super simple
for operators and admins to get into the usage of OpenEBS without much of learning
curve.

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
  - To run the mayaserver at **default bind address**:
    - sudo nohup mayaserver up &>mserver.log &
  - To run the mayaserver at a **particular bind address**:
    - sudo nohup mayaserver up -bind=172.28.128.4 &>mserver.log &

## Troubleshoot

- Check if mayaserver is running ?
  - Watch out for the process with 5656 as the port
  - `5656` is the default tcp port on which mayaserver's services are exposed

```bash
# use netstat command
$ netstat -tnlp

(Not all processes could be identified, non-owned process info
 will not be shown, you would have to be root to see it all.)
Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      -
tcp        0      0 127.0.0.1:5656          0.0.0.0:*               LISTEN      -
tcp6       0      0 :::22                   :::*                    LISTEN      -

# sudo will display the PID details
$ sudo netstat -tnlp

Active Internet connections (only servers)
Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      1258/sshd
tcp        0      0 127.0.0.1:5656          0.0.0.0:*               LISTEN      3078/mayaserver 
tcp6       0      0 :::22                   :::*                    LISTEN      1258/sshd

# use curl to check the services
curl http://$IP:5656/latest/meta-data/instance-id
```

## Licensing

Mayaserver is completely open source and bears an Apache license. Mayaserver's
core components and designs are a derivative of other open sourced libraries 
like Nomad and Kubernetes.
