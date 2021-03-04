# Deploy 2 servers - frontend and api

## Create the namespace **demo**

```sh
kubectl create namespace demo2server
```

## Create the api deployment from the yaml file
```sh
kubectl create -f api-deploy.yaml
```

## Create the api service from the yaml file
```sh
kubectl create -f api-service.yaml
```

## Create the frontend deployment from the yaml file
```sh
kubectl create -f frontend-deploy.yaml
```

## Create the frontend service from the yaml file
```sh
kubectl create -f frontend-service.yaml
```



## Check the results
```
kubectl get -n demo2server all

NAME                                  READY   STATUS    RESTARTS   AGE
pod/apiserver-59859d79c-w8mvd         1/1     Running   0          7m18s
pod/frontendserver-5445fcbc49-fkb8c   1/1     Running   0          2m2s

NAME                     TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
service/apiserver        ClusterIP   10.106.122.98   <none>        3000/TCP         6m58s
service/frontendserver   NodePort    10.110.124.90   <none>        8080:31162/TCP   19s

NAME                             READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/apiserver        1/1     1            1           7m18s
deployment.apps/frontendserver   1/1     1            1           2m2s

NAME                                        DESIRED   CURRENT   READY   AGE
replicaset.apps/apiserver-59859d79c         1         1         1       7m18s
replicaset.apps/frontendserver-5445fcbc49   1         1         1       2m2s

```

## Access the server from remote
```
curl localhost:31162

Host apiserver-59859d79c-w8mvd sent time: 04 Mar 21 10:09:03.34109 UTC%
```
>Attention: This works only on docker desktop

## Clean up
```
kubectl delete -n demo2server deployment,service apiserver frontendserver --force
```
