# Deploy the standalone server in declarative way

## Create the namespace **demo**

```sh
kubectl create namespace demo
```

## Create the deployment from the yaml file
```sh
kubectl create -f demodeploy.yaml
```

## Create the service from the yaml file
```sh
kubectl create -f demoservice.yaml
```

## Check the results
```
kubectl get -n demo all

NAME                           READY   STATUS    RESTARTS   AGE
pod/goserver-fbb99db68-bnvdj   1/1     Running   0          9m16s

NAME               TYPE       CLUSTER-IP    EXTERNAL-IP   PORT(S)          AGE
service/goserver   NodePort   10.97.57.97   <none>        8080:30377/TCP   12s

NAME                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/goserver   1/1     1            1           9m16s

NAME                                 DESIRED   CURRENT   READY   AGE
replicaset.apps/goserver-fbb99db68   1         1         1       9m16s
```

## Access the server from remote
```
curl localhost:30377
```
>Attention: This works only on docker desktop

## Clean up
```
kubectl delete -n demo deployment,service goserver --force
```
