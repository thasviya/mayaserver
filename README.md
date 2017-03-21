# Maya API Server

> Maya exposes its APIs here

A service exposing `Kubernetes` like APIs.

## Use Cases

### Serving EBS compatibility

Maya server adapts itself like Kubernetes server. This makes it super simple
for operators and admins to get into the usage of OpenEBS without much of learning
curve.

## Setting up Mayaserver locally

> These are some of the steps to start off with development & running of Mayaserver
in one's laptop. It assumes use of Linux as laptop's OS. In addition, the laptop
should have Virtual Box & Vagrant installed.

```bash
- git clone https://github.com/openebs/mayaserver.git
- cd to above cloned folder i.e mayaserver
  - vagrant up
  - vagrant ssh
- Inside the vagrant VM run below steps:
  - make init
  - make
  - make bin
  - sudo nohup mayaserver up -bind=172.28.128.4 &>mserver.log &
```

### Troubleshooting during local setup

```bash
- `make init` is a time taking operation
  - This downloads all the vendoring libraries
  - Typically required for the very first attempt only
  - In case of add/update of new/existing vendoring libraries:
    - use `make sync` than `make init`
```

## Mayaserver's REST APIs

- `NOTE: Use the bind address on which your Mayaserver is running`

- Get InstanceID

  ```bash
    # Metadata
    $ curl http://172.28.128.4:5656/latest/meta-data/instance-id
  ```

- Volume provisioning & deletion requires the presence of orchestrator config file
  - It is a .INI formatted file
  - It should be located at `/etc/mayaserver/orchprovider/`
  - There can be multiple .INI files based on the number of `regions` the orchestrator is deployed/running
    - `/etc/mayaserver/orchprovider/nomad_global.INI`
    - `/etc/mayaserver/orchprovider/nomad_us-east-1.INI`
  - Each region based .INI file can be categorised into multiple datacenters
  - Mayaserver works with following defaults:
    - `Nomad` as the default orchestrator &
    - `global` as Nomad's default region &
    - `dc1` as Nomad's default datacenter
  - Below is a sample .INI file that is valid for `Nomad` as mayaserver's orchestrator

  ```ini
  [datacenter "dc1"]
  address = http://172.28.128.3:4646
  cn-type = host
  cn-network-cidr = 172.28.128.1/24
  cn-interface = enp0s8
  ```

- Below is a sample volume spec that can be provisioned

  ```yaml
  # Similar to K8s' PersistentVolumeClaim
  kind: PersistentVolumeClaim
  apiVersion: v1
  metadata:
    name: minjvol
  spec:
    accessModes:
      - ReadWriteOnce
    resources:
      requests:
        storage: 3Gi
  ```

  - NOTE - Verbose specs can be found at:
    - `lib/mockit/verbose_openebs_pvc.yaml`
    - json formatted: `lib/mockit/verbose_openebs_pvc.json`
  - Both json & yaml specs are supported when http request's Content-Type is `application/yaml`

- Sample REST Calls
 
  ```bash
  
  # Provision
  
  $ curl -k -H "Content-Type: application/yaml" \
    -XPOST -d"$(cat lib/mockit/minimal_openebs_pvc.yaml)" \
    http://172.28.128.4:5656/latest/volumes/
    
  {
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
      "Reason": "complete"
    },
    "annotations": {
      "evalpriority": "50",
      "evaltype": "service",
      "evaltrigger": "job-register",
      "evaljob": "minjvol",
      "evalstatus": "complete",
      "evalstatusdesc": "",
      "evalblockedeval": "360ef0b3-e6dc-9353-d91c-d7c91c239102"
    },
    "creationTimestamp": null,
    "name": "minjvol"
  }

  # Info
  
  $ curl http://172.28.128.4:5656/latest/volume/info/myjivavol
  
  {
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
    "name": "minjvol"
  }

  # Delete
  
  $ curl http://172.28.128.4:5656/latest/volume/delete/myjivavol
  
  {
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
      "Reason": "complete"
    },
    "annotations": {
      "evaltype": "service",
      "evaltrigger": "job-deregister",
      "evaljob": "minjvol",
      "evalstatus": "complete",
      "evalstatusdesc": "",
      "evalblockedeval": "",
      "evalpriority": "50"
    },
    "creationTimestamp": null,
    "name": "minjvol"
  }

  # Info again
  
  $ curl http://172.28.128.4:5656/latest/volume/info/minjvol

  Unexpected response code: 404 (job not found)
    
  ```

## Troubleshooting

- Verify the presence of Mayaserver binary
  - which mayaserver
  - mayaserver -version

- Verify the presence of Mayaserver's orchestrator's .INI file(s)
  - i.e. /etc/mayaserver/orchprovider/nomad_global.INI
  - `global` is the name of the region

- Verify the contents of Mayaserver's orchestrator's .INI file

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

- How to get the iscsi portal & iqn information of a iscsi based jiva volume ?
  - Info based REST API will fetch these information.
  - However, these will **not** be fetched in case of:
    - `error` or
    - `in-progress` / `pending` status

## Licensing

Mayaserver is completely open source and bears an Apache license. Mayaserver's
core components and designs are a derivative of other open sourced libraries 
like Nomad and Kubernetes.
