Build Docker Image:
-> docker build -t image_name:latest . 
-> Before running this command DB_DRIVERockerfile should be setup properly.
-> If you need light weight docker image size the use multi stage Dockerfile stages can be Build Stage, Run Stage.
-> If you want to run the container on specific network use --network flag

To list all images:
-> docker images

To Remove an docker image:
-> docker rmi image_name

To list running pods:
-> docker ps

To list all pods:
-> docker ps -a

To Run the Container from build Image:
-> docker run --name container_name -p port:port imagename:tag
-> Use -e flag to pass the environment variable.

To Remove a container:
-> docker rm container_name

To check Network Settings of a Container:
-> docker container inspect container_name

To list all networks:
-> docker network ls

To check more details of a network:
-> docker network inspect network_name
-> Containers running on the same network discovers each other via name.
-> Docker provides several commands to manage the network.

To create a new network:
-> docker network create network_name

To connect to a network:
-> docker network connect network_name container_name

To Run docker compose file:
-> In CLI type, docker compose up
-> To Remove Containers created using docker compose up use "docker compose down"
-> To monitor logs use "docker compose logs"
-> To list all services with their status use "docker compose ps"

When CMD command is used with entrypoint command in DOckerfile then cmd will act as a param passed to entrypoint.

Kubernetes Notes:
-> Kubernetes is an open source container orchestration engine for automating deployment, scaling and management of containerized 
   applications.

-> Kubernetes Components:
    -> Worker Node
        -> Runs the containerized application.
        -> Each node contains a kubelet agent which makes all container runs inside pods.
        -> Kubernetes supports several container runtimes such as docker, containerd or CRI-O
        -> Kube-proxy maintains the network rules and allows communication with pods.
    -> Master Node
        -> Second part is the control plane which runs on master node.
        -> it manages the worker nodes and the pods of the cluster.
        -> Control plane contains several components such as
            -> API Server: which is the front end of the control plane. It exposes Kubernetes api to interact with all other components 
               of cluster.
            -> etcd backup of the api Server.
            -> scheduler watches for the newly created pods with no assigned nodes and select the nodes for them to run on.
            -> Control Manager is combination of several controllers such as
                -> Node Controller responsible for noticing and responding when nodes go down.
                -> Job Controller which watches for jobs jobs or one-off tasks then creates pods to run them.
                -> Endpoint Controller which joins services and pods
                -> Service Account & Token Controller which creates default account and API access token for new namespaces.
            -> Cloud Controller Manager when services deployed on cloud which links cluster into the cloud provider's API
                -> Node Controller
                -> Route Controller
                -> Service Controller



            




