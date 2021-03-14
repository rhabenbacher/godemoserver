# Run the goserver in kubernetes

## Create a deployment

```sh
kubectl create deployment goserver --image=goserver:0.2
```

## Create a service and expose it on port 8080
```sh
kubectl expose deployment goserver --port=8080 --target-port=9000 --type=NodePort
```

## Find the nodeport to access the server from remote
```
kubectl get service

NAME                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
service/goserver     NodePort    10.105.253.217   <none>        8080:31486/TCP   15s
```
The port is **31486**

## Access the server
```
% curl localhost:31486

Hello, I'm serving in standalone mode !!!

I'm running in a Pod with name goserver-5f4b44d46b-f8zlz 

04 Mar 21 09:03 UTC
```

> Attention: This works only on Docker Desktop. In other installations you can access the nodeport via the external IP address of one node. 