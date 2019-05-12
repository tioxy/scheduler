kube-stack
=========

Installs [Kubectl](https://github.com/kubernetes/kubectl), [Kubelet](https://github.com/kubernetes/kubelet) and [Kubeadm](https://github.com/kubernetes/kubeadm).

Requirements
------------

- Debian 9

Role Variables
--------------

Stick with */defaults/main.yml* values.
```yaml
kube_gpg_key: kubernetes gpg key
kube_repository: kubernetes deb package repository
kube_requirements: list of requirements to install kubernetes packages
kube_packages: list of kubernetes packages to install
```

Example Playbook
----------------

```yaml
- hosts: all

  roles:
    - role: kube-stack
```

License
-------

MIT

Author Information
------------------

Created by [@tioxy](https://github.com/tioxy)
