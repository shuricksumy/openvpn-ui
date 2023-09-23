# OpenVPN WEB UI

- There is the original [README.md](https://github.com/shuricksumy/openvpn-ui/blob/my_main/README_ORIGINAL.md) file
- The project was cloned from [https://github.com/d3vilh/openvpn-ui](https://github.com/d3vilh/openvpn-ui)  - big thanks for a great job
- The Docker builder with the Server part is here [**OPENVPN-SERVER-DOCKER**](https://github.com/shuricksumy/openvpn-server-docker)
- The base example for docker-compose is here [***DOCKER-COMPOSE***](https://github.com/shuricksumy/openvpn-server-docker/tree/main/examples/basic)

![Status page](screenshots/2.png?raw=true)

## Updates
#### September 2023 (v 2.2)
- Small improvements 
#### September 2023 (v 2.1)
- Added md5 sum checker for client config files to be sure that current config is used or not
- Fixed small issue
#### September 2023 (v 2.0)
- Now is possible to organize routing between devices on Web UI
    - Added Client details page with Static IP, Routes, Subnet settings, Default Route
    - Added script to generate client config files based on these settings
- It's possible un-revoke certificate
- Redesigned a bit UI

#### Summer 2023 (v 1.0)

- updated all config files and scripts to use `/etc/openvpn/easy-rsa` path
- added the script from [openvpn-install](https://github.com/shuricksumy/openvpn-install) as the main script for generating new clients
- added UI improvements:
  - now user can configure `server.conf` and `client-template` files as plain text
  - new table with certificates
  - add a modal window to edit each client config file separately
  - improved visual part of the log viewer
  - updated client generation and .ovpn file generation
  - added confirmation to Revoke or Delete clients
  - added 4 tabs for Application, Server, Cliemt config and System utils
  - added backuping/downloading of all ovpn directory
- added Docker env variables and improved run-script:
  - disabled auto-provisioning of OpenVPN server part - now wait for a readily configured server
  - added env vars:

  ```bash
  SITE_NAME='Server 1' # The name of the server - displayed on UI. Default value "Admin"
  OPENVPN_SERVER_DOCKER_NAME="vpnserver1" # The name of the Docker container to restart
  OPENVPN_MANAGEMENT_ADDRESS="IP:PORT" # The preconfigured address to connect OpenVPN manager
  ```
  
## Example docker-compose file

### It's only UI part - full configuration will be here soon [TODO]

```docker
version: '3'

networks:
    default:
        driver: bridge
    npm_proxy:
        name: npm_proxy
        driver: bridge
        ipam:
            config:
                - subnet: 172.18.0.0/24
services:
    gui:
        image: shuricksumy/openvpn-ui
        container_name: openvpn-ui
        working_dir: /etc/openvpn/easy-rsa
        environment:
            - OPENVPN_ADMIN_USERNAME=admin # Leave this default as-is and update on first-run
            - OPENVPN_ADMIN_PASSWORD=admin # Leave this default as-is and update on first-run
            - SITE_NAME=UDP Server
            - OPENVPN_SERVER_DOCKER_NAME=openvpn-server-1
            - OPENVPN_MANAGEMENT_ADDRESS=172.18.0.1:2080
        ports:
            - "8080:8080/tcp"
        restart: always
        networks:
            npm_proxy:
                ipv4_address: 172.18.0.12
        volumes:
         - /var/run/docker.sock:/var/run/docker.sock
         - ./openvpn/openvpn1:/etc/openvpn
         - ./openvpn/easy-rsa:/etc/openvpn/easy-rsa
         - ./openvpn/openvpn1/db:/opt/openvpn-gui-tap/db
```

## Screenshots

![Status page](screenshots/2.png?raw=true)
![Status page](screenshots/3.png?raw=true)
![Status page](screenshots/4.png?raw=true)
![Status page](screenshots/5.png?raw=true)
![Status page](screenshots/6.png?raw=true)
![Status page](screenshots/7.png?raw=true)
![Status page](screenshots/8.png?raw=true)
![Status page](screenshots/9.png?raw=true)
![Status page](screenshots/10.png?raw=true)
