# OpenVPN WEB UI

![Status page](screenshots/1.png?raw=true)

## Updates

#### December 2023 (v4.0)
> **Warning**
 There is no back compatibility with previous versions. Need to recreate service from scratch.
 
- Now it's a solid solution: OpenVPN server is included to docker
- UI can see the OpenVPN server status and restart it
- For better UX, Wizard was added to configure the OpenVPN server the first time
- Add Clients as first-level entities, stored in DB
- Certificates now can be generated only for created Clients
- Add Routes system management to provide each client with Route rules, stored in DB
- Refactored code a bit
- Redesigned UI a bit

<details>

<summary>Previous versions details</summary>

#### September 2023 (v3.0)
- New UI web components
- UI updates
#### September 2023 (v2.0 - v 2.4)
- Fixed some issues
- Add script based on go for client's file generation
- Small improvements 
- Added md5 sum checker for client config files to be sure that current config is used or not
- Fixed small issue
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
  
</details>

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
        image: shuricksumy/openvpn-ui:latest
        container_name: openvpn-ui
        working_dir: /etc/openvpn/easy-rsa
        environment:
            - OPENVPN_ADMIN_USERNAME=admin # Leave this default as-is and update on first-run
            - OPENVPN_ADMIN_PASSWORD=admin # Leave this default as-is and update on first-run
            - SITE_NAME=Admin
        ports:
            - "8080:8080/tcp"
        restart: always
        networks:
            npm_proxy:
                ipv4_address: 172.18.0.10
        devices:
            - /dev/net/tun
        cap_add:
            - NET_ADMIN
        volumes:
         -  /var/run/docker.sock:/var/run/docker.sock
        #  - ./openvpn/openvpn:/etc/openvpn
        #  - ./openvpn/easy-rsa:/etc/openvpn/easy-rsa
        #  - ./openvpn/openvpn/db:/opt/openvpn-gui/db
```

References:
- The project is originally based on [https://github.com/d3vilh/openvpn-ui](https://github.com/d3vilh/openvpn-ui)  - big thanks for a great job
- The bash script for setup OpenVPN is based on [https://github.com/angristan/openvpn-install](https://github.com/angristan/openvpn-install) - big thanks for a great job

## Screenshots

![Status page](screenshots/0.png?raw=true)
![Status page](screenshots/1.png?raw=true)
![Status page](screenshots/2.png?raw=true)
![Status page](screenshots/3.png?raw=true)
![Status page](screenshots/4.png?raw=true)
![Status page](screenshots/5.png?raw=true)
![Status page](screenshots/6.png?raw=true)
![Status page](screenshots/7.png?raw=true)
![Status page](screenshots/8.png?raw=true)
![Status page](screenshots/9.png?raw=true)
![Status page](screenshots/10.png?raw=true)
![Status page](screenshots/11.png?raw=true)
