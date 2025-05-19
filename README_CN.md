<h1 align="center">
  <br>
  <img src="https://github.com/Safe3/openresty-manager/blob/main/logo.png" alt="OpenResty Manager" width="70px">
</h1>
<h4 align="center">OpenResty 管理器</h4>

<p align="center">
<a href="https://github.com/Safe3/openresty-manager/releases"><img src="https://img.shields.io/github/downloads/Safe3/openresty-manager/total">
<a href="https://github.com/Safe3/openresty-manager/graphs/contributors"><img src="https://img.shields.io/github/contributors-anon/Safe3/openresty-manager">
<a href="https://github.com/Safe3/openresty-manager/releases/"><img src="https://img.shields.io/github/release/Safe3/openresty-manager">
<a href="https://github.com/Safe3/openresty-manager/issues"><img src="https://img.shields.io/github/issues-raw/Safe3/openresty-manager">
<a href="https://github.com/Safe3/openresty-manager/discussions"><img src="https://img.shields.io/github/discussions/Safe3/openresty-manager">
</p>
<p align="center">
  <a href="#dart-特色">特色</a> •
  <a href="#rocket-使用">使用</a> •
  <a href="#gift_heart-感谢">感谢</a> •
  <a href="#kissing_heart-联系">联系</a> •
  <a href="#key-授权">授权</a>
</p>

<p align="center">
  <a href="https://github.com/Safe3/openresty-manager/blob/main/README.md">English</a>
  <br/><br/>
  ⭐请帮我们点个star以支持我们不断改进，谢谢！
</p>

---

最简单易用、功能强大、漂亮的OpenResty管理器(Nginx增强版)，OpenResty Edge 的开源替代品，它可以让您轻松地安全反向代理在家或互联网上运行的网站，包括访问控制、拒绝服务攻击防护、自动申请并续期免费的SSL证书，而无需对OpenResty或Let's Encrypt了解太多。

<h3 align="center">
  <img src="https://github.com/Safe3/openresty-manager/blob/main/openresty-manager_cn.png" alt="OpenResty Manager" width="700px">
  <br>
</h3>

## :dart: 特色
:green_circle: 提供了美观、强大且易于使用的web管理UI

 :purple_circle: 支持通过HTTP-01和DNS-01挑战申请免费SSL证书并自动续期

 :yellow_circle: 无需了解OpenResty，即可轻松为您的网站创建反向代理

 :orange_circle: 提供了多种强大安全功能，如访问控制、拒绝服务攻击防护等

 :red_circle: 高级OpenResty配置可供超级用户使用

 :large_blue_circle: 支持并继承OpenResty超越Nginx的所有强大功能



## :rocket: 使用

OpenResty Manager不仅易于使用，而且易于安装，支持主机和容器环境。

- ### 云服务部署

&nbsp;&nbsp;<a href="https://app.rainyun.com/apps/rca/store/6202?ref=689306" target="_blank"><img height="42" src="https://rainyun-apps.cn-nb1.rains3.com/materials/deploy-on-rainyun-cn.svg" alt="雨云"></a>

&nbsp;&nbsp;<a href="https://8465.cn/aff/NCKQREHC" target="_blank"><img height="32" src="https://8465.cn/themes/web/www/upload/local665305c838bfb.png" alt="蓝谷科技"></a>

- ### 主机版

> :biohazard: ***如果服务器正在使用云服务，请记住开放OpenResty Manager所需的TCP端口80、443和34567***

一键安装：自动安装可以在几分钟内完成。

```bash
sudo bash -c "$(curl -fsSL https://om.uusec.com/installer_cn.sh)"
```



- ### 容器版

一键安装：自动安装可以在几分钟内完成。

```bash
sudo bash -c "$(curl -fsSL https://om.uusec.com/docker_installer_cn.sh)"
```

随后可以通过 bash /opt/om/om.sh 命令管理OpenResty Manager容器，包括启动、停止、更新、卸载等。



- ### 快速入门

1. 登录管理：访问http://ip:34567 ，默认用户名为“admin”，默认密码为“#Passw0rd”。
2. 添加SSL证书：转到证书管理菜单，申请Let's Encrypt免费SSL证书或上传现有证书。
3. 添加上游：转到上游管理菜单，为原始站点添加上游负载均衡。
4. 添加站点：进入站点菜单，点击“新建站点”按钮，按照提示添加反向代理的站点域名。
5. 测试连接：将您的域名DNS A或CNAME记录更改为OpenResty Manager的服务器IP，访问您的网站查看是否可以打开。




- ### 卸载

一键卸载：自动卸载可以在几分钟内完成。

```bash
sudo bash -c "$(curl -fsSL https://om.uusec.com/uninstaller.sh)"
```




## :gift_heart: 感谢

感谢所有出色的[社区贡献者发送PR](https://github.com/Safe3/openresty-manager/graphs/contributors)并保持此项目的更新。❤️

如果你有一个想法或某种改进，欢迎你贡献和参与这个项目，随时发送你的PR。

<p align="center">
<a href="https://github.com/Safe3/openresty-manager/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=Safe3/openresty-manager&max=500">
</a>
</p>
捐赠请扫描如下二维码：
<img src="https://waf.uusec.com/_media/sponsor.jpg" alt="捐赠"  height="300px" />



## :kissing_heart: 联系

如果您想支持更多功能，如Web应用程序防火墙，请访问[南墙](https://waf.uusec.com/)项目。

- 官方 QQ 群：11500614

- 官方微信群：微信扫描以下二维码加入

  <img src="https://waf.uusec.com/_media/weixin.jpg" alt="微信群"  height="200px" />


## :key: 授权

OpenResty Manager遵循GPL许可证，每个人都可以免费使用！

