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
    - `NOTE: This is a time taking operation`
    - This downloads all the vendoring libraries
    - Typically required for the very first attempt only
    - In case of add/update of new/existing vendoring libraries:
      - use `make sync` than `make init`
  - make
  - make bin
  - To run the mayaserver at a **particular bind address**:
    - sudo nohup mayaserver up -bind=172.28.128.4 &>mserver.log &

## Mayaserver's REST APIs

- `NOTE: Use the bind address on which your Mayaserver is running`

- Get InstanceID

  ```bash
    # Metadata
    $ curl http://172.28.128.4:5656/latest/meta-data/instance-id
  ```

- Volume provisioning & deletion requires the presence of a .INI file
  - This orchestrator file provides the coordinates of Nomad server/cluster
  - i.e. `/etc/mayaserver/orchprovider/nomad_global.INI`

- Below is a sample volume spec that can be provisioned

  ```yaml
  # Similar to K8s' PersistentVolumeClaim
  kind: PersistentVolumeClaim
  apiVersion: v1
  metadata:
    name: ssdvol
    labels:
      region: global
      datacenter: dc1
      jivafeversion: openebs/jiva:latest
      jivafenetwork: host
      jivafeip: 172.28.128.101
      jivabeip: 172.28.128.102
      jivafesubnet: 24
      jivafeinterface: enp0s8
    annotations:
      volume.beta.openebs.io/orchestrator-class: nomad
  spec:
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: 3Gi
  ```

- Sample REST Calls
 
  ```bash
  
  # Provision
  
  $ curl -k -H "Content-Type: application/yaml" \
    -XPOST -d"$(cat lib/mockit/sample_openebs_pvc.yaml)" \
    http://172.28.128.4:5656/latest/volumes/
    
  {
    "Allocs": null,
    "Evals": [
      {
        "BlockedEval": "",
        "CreateIndex": 1992,
        "FailedTGAllocs": null,
        "ID": "6ed438d2-    d7fe-7e51-91cc-034018405db4",
        "JobID": "ssdvol",
        "JobModifyIndex": 1991,
        "ModifyIndex": 1992,
        "NextEval": "",
        "NodeID": "",
        "NodeModifyIndex": 0,
        "PreviousEval": "",
        "Priority": 50,
        "QueuedAllocations": null,
        "Status": "pending",
        "StatusDescription": "",
        "TriggeredBy": "job-register",
        "Type": "service",
        "Wait": 0
      }
    ],
    "Spec": {
      "AccessModes": null,
      "Capacity": null,
      "ClaimRef": null,
      "OpenEBS": {
        "volumeID": ""
      },
      "PersistentVolumeReclaimPolicy": "",
      "StorageClassName": ""
    },
    "Status": {
      "Message": "",
      "Phase": "",
      "Reason": ""
    },
    "creationTimestamp": null,
    "name": "ssdvol"
  }


  # Info
  
  $ curl http://172.28.128.4:5656/latest/volume/info/ssdvol
  
  {
    "Allocs": null,
    "Evals": null,
    "Spec": {
      "AccessModes": null,
      "Capacity": null,
      "ClaimRef": null,
      "OpenEBS": {
        "volumeID": ""
      },
      "PersistentVolumeReclaimPolicy": "",
      "StorageClassName": ""
    },
    "Status": {
      "Message": "",
      "Phase": "",
      "Reason": "pending"
    },
    "creationTimestamp": null,
    "name": "ssdvol"
  }


  # Delete
  
  $ curl http://172.28.128.4:5656/latest/volume/delete/ssdvol
  
  {
    "Allocs": null,
    "Evals": [
      {
        "BlockedEval": "",
        "CreateIndex": 2023,
        "FailedTGAllocs": null,
        "ID": "27d49622-3985-8fc9-08ac-1ed6cbd3eb3b",
        "JobID": "ssdvol",
        "JobModifyIndex": 2022,
        "ModifyIndex": 2023,
        "NextEval": "",
        "NodeID": "",
        "NodeModifyIndex": 0,
        "PreviousEval": "",
        "Priority": 50,
        "QueuedAllocations": null,
        "Status": "pending",
        "StatusDescription": "",
        "TriggeredBy": "job-deregister",
        "Type": "service",
        "Wait": 0
      }
    ],
    "Spec": {
      "AccessModes": null,
      "Capacity": null,
      "ClaimRef": null,
      "OpenEBS": {
        "volumeID": ""
      },
      "PersistentVolumeReclaimPolicy": "",
      "StorageClassName": ""
    },
    "Status": {
      "Message": "",
      "Phase": "",
      "Reason": ""
    },
    "creationTimestamp": null,
    "name": "ssdvol"
  }

  # Info again
  
  $ curl http://172.28.128.4:5656/latest/volume/info/ssdvol

  Unexpected response code: 404 (job not found)
    
  ```

## Troubleshoot

- Verify the presence of Mayaserver binary
  - which mayaserver
  - mayaserver -version

- Verify the presence of Mayaserver's orchestrator's .INI file(s)
  - i.e. /etc/mayaserver/orchprovider/nomad_global.INI
  - `global` is the name of the region

- Verify the contents of Mayaserver's orchestrator's .INI file
  - Below is a sample .INI file that is valid for Nomad as mayaserver's orchestrator

  ```ini
  [datacenter "dc1"]
  address = http://172.28.128.3:4646
  ```

- Verify if Mayaserver is running as a process
    - Watch out for the process with 5656 as the port
    - `5656` is the default tcp port on which Mayaserver's services are exposed

  ```bash
  # Use netstat command
  $ netstat -tnlp

  (Not all processes could be identified, non-owned process info
   will not be shown, you would have to be root to see it all.)
  Active Internet connections (only servers)
  Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
  tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      -
  tcp        0      0 127.0.0.1:5656          0.0.0.0:*               LISTEN      -
  tcp6       0      0 :::22                   :::*                    LISTEN      -

  # Using sudo will display the PID details
  $ sudo netstat -tnlp

  Active Internet connections (only servers)
  Proto Recv-Q Send-Q Local Address           Foreign Address         State       PID/Program name
  tcp        0      0 0.0.0.0:22              0.0.0.0:*               LISTEN      1258/sshd
  tcp        0      0 127.0.0.1:5656          0.0.0.0:*               LISTEN      3078/mayaserver 
  tcp6       0      0 :::22                   :::*                    LISTEN      1258/sshd
  ```


## Licensing

Mayaserver is completely open source and bears an Apache license. Mayaserver's
core components and designs are a derivative of other open sourced libraries 
like Nomad and Kubernetes.
