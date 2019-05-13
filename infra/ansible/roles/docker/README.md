docker
=========

Installs [Docker](https://www.docker.com/what-docker).

Requirements
------------

- Debian 9

Role Variables
--------------

Stick with */defaults/main.yml* values.
```yaml
docker_gpg_key: docker gpg key
docker_repository: docker deb package repository
docker_requirements: list of requirements to install Docker
```

Example Playbook
----------------

```yaml
- hosts: all

  roles:
    - role: docker
```

License
-------

MIT

Author Information
------------------

Created by [@tioxy](https://github.com/tioxy)
