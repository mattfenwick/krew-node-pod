# node-pod kubectl

A `kubectl` plugin to show the assignment of pods to nodes.

## Quick Start

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

     - easy way: `kubectl krew install --manifest=./deploy/krew/plugin.yaml`

     - hard way:
        1. download a binary from [the github project release page](https://github.com/mattfenwick/krew-node-pod/releases/tag/v0.0.3)

        2. run a `krew install` against the downloaded binary

            ```
            kubectl krew install --manifest=./deploy/krew/plugin.yaml --archive=/Users/mfenwick/Downloads/node-pod_darwin_amd64.tar.gz
            ```

2. test it

    ```
    kubectl node-pod
    ```

3. clean up

    ```
    kubectl krew uninstall node-pod
    ```