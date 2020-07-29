# kubectl node-pod

A `kubectl` plugin to show the assignment of pods to nodes.

```
kubectl node-pod --format table --status
+----------------------+--------------------+----------------------------------------------+---------+
|         NODE         |     NAMESPACE      |                   POD NAME                   | STATUS  |
+----------------------+--------------------+----------------------------------------------+---------+
| calico-control-plane |                    |                                              | Ready   |
|                      | kube-system        | calico-node-qfcqs                            | Running |
|                      | kube-system        | etcd-calico-control-plane                    | Running |
|                      | kube-system        | kube-apiserver-calico-control-plane          | Running |
|                      | kube-system        | kube-controller-manager-calico-control-plane | Running |
|                      | kube-system        | kube-proxy-9wn87                             | Running |
|                      | kube-system        | kube-scheduler-calico-control-plane          | Running |
| calico-worker        |                    |                                              | Ready   |
|                      | kube-system        | calico-kube-controllers-857b8b787f-wpj29     | Running |
|                      | kube-system        | calico-node-tpphd                            | Running |
|                      | kube-system        | coredns-66bff467f8-2nwsr                     | Running |
|                      | kube-system        | kube-proxy-fphtj                             | Running |
|                      | sonobuoy-3         | sonobuoy                                     | Running |
|                      | sonobuoy-4         | sonobuoy                                     | Running |
| calico-worker2       |                    |                                              | Ready   |
|                      | kube-system        | calico-node-qvdqd                            | Running |
|                      | kube-system        | coredns-66bff467f8-nnggj                     | Running |
|                      | kube-system        | kube-proxy-bzcwf                             | Running |
|                      | local-path-storage | local-path-provisioner-bd4bb6b75-dcr9g       | Running |
|                      | sonobuoy-5         | sonobuoy                                     | Running |
+----------------------+--------------------+----------------------------------------------+---------+
```

## Quick Start

1. Download the [latest binary for your OS](https://github.com/mattfenwick/krew-node-pod/releases)

2. unzip the archive

3. move the executable somewhere in your path

    ```bash
    # OS X example;  will be different on other OSs
    mv ~/Downloads/kubectl-node_pod_darwin_amd64/kubectl-node_pod /usr/local/bin
    ```

4. run it!

    ```
    kubectl node-pod
    ```

or TODO add to krew:

```
kubectl krew install node-pod
kubectl node-pod
```

# How to release

See https://goreleaser.com/quick-start/

Set up a github token and run:

```
VERSION=v9.9.9
git tag $VERSION -m "something about the version"
git push --tags

export GITHUB_TOKEN=...

goreleaser
```

In [the plugin.yaml file](./deploy/krew/plugin.yaml):
 - manually update the sha256's using [the checksums file](https://github.com/mattfenwick/krew-node-pod/releases/download/v0.0.3/node-pod_0.0.3_checksums.txt)
   - **TODO** is there a better way to do this?
 - manually update the versions using [the checksums file](https://github.com/mattfenwick/krew-node-pod/releases/download/v0.0.3/node-pod_0.0.3_checksums.txt)
   - **TODO** is there a better way to do this?

# How to test a local release

1. Choose one of the following:

     - easy way: `kubectl krew install --manifest=./deploy/krew/node-pod.yaml`

     - hard way:
        1. download a binary from [the github project release page](https://github.com/mattfenwick/krew-node-pod/releases/tag/v0.0.3)

        2. run a `krew install` against the downloaded binary

            ```
            kubectl krew install --manifest=./deploy/krew/node-pod.yaml --archive=/Users/mfenwick/Downloads/node-pod_darwin_amd64.tar.gz
            ```

2. test it

    ```
    kubectl node-pod
    ```

3. clean up

    ```
    kubectl krew uninstall node-pod
    ```