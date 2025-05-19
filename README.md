<h1 align="center">
  <br>
  <img src="https://github.com/Safe3/openresty-manager/blob/main/logo.png" alt="OpenResty Manager" width="70px">
</h1>
<h4 align="center">OpenResty Manager</h4>

<p align="center">
<a href="https://github.com/Safe3/openresty-manager/releases"><img src="https://img.shields.io/github/downloads/Safe3/openresty-manager/total">
<a href="https://github.com/Safe3/openresty-manager/graphs/contributors"><img src="https://img.shields.io/github/contributors-anon/Safe3/openresty-manager">
<a href="https://github.com/Safe3/openresty-manager/releases/"><img src="https://img.shields.io/github/release/Safe3/openresty-manager">
<a href="https://github.com/Safe3/openresty-manager/issues"><img src="https://img.shields.io/github/issues-raw/Safe3/openresty-manager">
<a href="https://github.com/Safe3/openresty-manager/discussions"><img src="https://img.shields.io/github/discussions/Safe3/openresty-manager">
</p>
<p align="center">
  <a href="#dart-features">Features</a> •
  <a href="#rocket-usage">Usage</a> •
  <a href="#gift_heart-credits">Credits</a> •
  <a href="#kissing_heart-contact">Contact</a> •
  <a href="#key-license">License</a>
</p>






<p align="center">
  <a href="https://github.com/Safe3/openresty-manager/blob/main/README_CN.md">中文</a>
  <br/><br/>
  ⭐Please help us with a star to support our continuous improvement, thank you!
</p>




---

The easiest using, powerful and beautiful OpenResty Manager (Nginx Enhanced Version) , open source alternative to OpenResty Edge, which can enable you to easily reverse proxy your websites with security running at home or internet, including Access Control, HTTP Flood Protection, Free SSL, without having to know too much about OpenResty or Let's Encrypt.

<h3 align="center">
  <img src="https://github.com/Safe3/openresty-manager/blob/main/openresty-manager.png" alt="OpenResty Manager" width="700px">
  <br>
</h3>


## :dart: Features
:green_circle: Provide a beautiful, powerful, and easy-to-use web management UI

 :purple_circle: Free SSL support both for HTTP-01 and DNS-01 challenge or provide your own SSL certificates

 :yellow_circle: Easily create reverse proxy for your websites without knowing anything about OpenResty

 :orange_circle: Provide some powerful security function, such as Access Control, HTTP Flood Protection, etc

 :red_circle: Advanced OpenResty configuration available for super users

 :large_blue_circle: Support and inherit all powerful features of OpenResty



## :rocket: Usage

OpenResty Manager is not only easy to use but also easy to install, supports both host and container environments.

- ### Cloud Deploy

&nbsp;&nbsp;<a href="https://app.rainyun.com/apps/rca/store/6202?ref=689306" target="_blank"><img height="42" src="https://rainyun-apps.cn-nb1.rains3.com/materials/deploy-on-rainyun-en.svg" alt="RainYun"></a>

&nbsp;&nbsp;<a href="https://8465.cn/aff/NCKQREHC" target="_blank"><img height="32" src="https://8465.cn/themes/web/www/upload/local665305c838bfb.png" alt="蓝谷科技"></a>

- ### Host Version

> :biohazard: ***If the server is using cloud services, remember to open the TCP port 80, 443 and 34567 required for OpenResty Manager***

One click installation: Automatic installation can be completed in minutes.

```bash
sudo bash -c "$(curl -fsSL https://om.uusec.com/installer.sh)"
```

- ### Docker Version

One click installation: Automatic installation can be completed in minutes.

```bash
sudo bash -c "$(curl -fsSL https://om.uusec.com/docker_installer.sh)"
```

Subsequently, `bash /opt/om/om.sh` is used to manage the OpenResty Manager container, including starting, stopping, updating, uninstalling, etc.

- ### Quick Start

1. Login to the management: Access http://ip:34567 , the default username is "admin", and the default password is "#Passw0rd".
2. Add SSL certificates: Go to the certificates management menu, apply for a Let's Encrypt free SSL certificate or upload an existing certificate.
3. Add upstreams: Go to the upstreams management menu, add a load balancing upstream that for your original sites.
4. Add a site: Go to the sites menu, click the "New site" button, and follow the prompts to add the site domain names for reverse proxy.
5. Test connectivity: Change your domain dns A or CNAME record to the OpenResty Manager server IP, visit your website to see if it can be opened.

- ### Uninstall

One click uninstallation: Automatic uninstallation can be completed in minutes.

```bash
sudo bash -c "$(curl -fsSL https://om.uusec.com/uninstaller.sh)"
```

## :gift_heart: Credits

Thanks to all the amazing [community contributors for sending PRs](https://github.com/Safe3/openresty-manager/graphs/contributors) and keeping this project updated. ❤️

If you have an idea or some kind of improvement, you are welcome to contribute and participate in the Project, feel free to send your PR.

<p align="center">
<a href="https://github.com/Safe3/openresty-manager/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Safe3/openresty-manager&max=500">
</a>
</p>

## :kissing_heart: Contact

If you want to support more features such as Web Application Firewall, please visits [UUSEC WAF](https://uuwaf.uusec.com/).

## :key: License

OpenResty Manager is under GPL license, everyone can use it for free！

