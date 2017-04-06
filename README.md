# Kubernetes Webhook Token Authenticator for GitHub

This project implements a Kubernetes [Webhook Token
Authenticator](https://kubernetes.io/docs/admin/authentication/#webhook-token-authentication)
for authenticating users using GitHub Personal Access Token.

When user
tries to authenticate to the Kubernetes API, the Kubernetes apiserver
calls this authenticator to verify the bearer token. This authenticator checks
if the access token is valid using GitHub API and returns the GitHub username
to apiserver.

You should configure Kubernetes apiserver with an [authorization
plugin](https://kubernetes.io/docs/admin/authorization/) to control what
Kubernetes resources can a user access.

## How to use

First of all, you need to run the authenticator using the example [DaemonSet
manifest](manifests/github-authn.yaml). It is recommended to run the
authenticator on your Kubernetes master using host networking so that the
apiserver can access the authenticator through the loopback interface.

```
kubectl create -f https://github.com/oursky/kubernetes-github-authn/blob/master/manifests/github-authn.yaml
```

Confirm that the authenticator is running:

```
kubectl get pod -l k8s-app=github-authn
```

Next, configure apiserver to verify bearer token using this authenticator.
There are two configuration options you need to set:

* `--authentication-token-webhook-config-file` a kubeconfig file describing how to
  access the remote webhook service.
* `--authentication-token-webhook-cache-ttl` how long to cache authentication
  decisions. Defaults to two minutes.

Check the [example config file](manifests/token-webhook-config.json) and save
this file in the Kubernetes master. Set the path to this config file
with configurion option above.

It is recommended you read the [Kubernetes
documentation](https://kubernetes.io/docs/admin/authentication/#webhook-token-authentication) for how to configure
webhook token authentication.

## Authorization with role-based access control (RBAC)

Kubernetes support multiple [authorization
plugins](https://kubernetes.io/docs/admin/authorization) and we recommend
you choose role-based access control (RBAC) because permission settings can be
set using the Kubernetes API. Permission is granted on which roles that the
authenticated user has.

Suppose that we have a user called `johndoe` and this user has administrative
access to the project `project1`. First of all, we need to define a new role
called `admin` which can control all resources.

```
kubectl create -f https://github.com/oursky/kubernetes-github-authn/blob/master/manifests/admin-cluster-role.yaml
```

We need to assign `johndoe` to this `admin` role so that he has control to
all the resources in the namespace `project1`.

```
kubectl create namespace project1
kubectl create rolebinding johndoe-admin-binding --clusterrole=admin --user=johndoe --namespace=project1
```

If we want to assign `johndoe` to the `admin` role in all namespaces instead of
just the `project1` namespace, create a `ClusterRoleBinding` instead of
a `RoleBinding`:

```
kubectl create clusterrolebinding johndoe-admin-binding --clusterrole=admin --user=johndoe
```

Read the [Kubernetes
documentation](https://kubernetes.io/docs/admin/authorization/rbac/) to learn
more about how to configure your apiserver to use RBAC.
