# Golang Datastore Microservices

Application contains two microservices: Datastore Service and Encryption-Server Services.
A composite application can be used by command-line tool which is also part of the project.
Services are able to communication to each other using Discovery Service that ensures that
services can know only symbolic names, real URLs are stored in ETCD database. Every Service
that wants to be part of the microservices environment must register itself using provided
library.

### Run Application
- install ETCD 
    - LINUX - UBUNTU: `sudo apt-get install etcd`
     
- run ETCD
    - LINUX - UBUNTU: `/usr/bin/etcd`
    - etcd can be accessible on localhost:4001 by default
        
- build golang code
    - `./build-apps.sh`
    
- run services
    - `./datastore/datastore`
    - `./encryption/encryption`
    
- encrypt and store value in datastore
    - `./yoti store my-key my-value` -> returns encrtyption key `my-encryption-key`
    - `./yoti retrieve my-key my-encryption-key` -> return original value
    
### Enhancements
- Client is able use only first registered instance -> implement Round-Robin mechanism on client
- Better exception handling and react to returned status codes
- Health checking is implemented by push the service itself, then the service is removed from Discovery Service but the service is not restarted
- And many more other features that are implemented in Mesos, Kubernetes and others :)