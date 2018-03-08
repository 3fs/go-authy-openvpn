# Go Authy OpenVPN

With Go Authy OpenVPN plugin you can add two-factor authentication to your VPN server in just minutes.
It is a replacement for official [Authy OpenVPN plugin](https://github.com/authy/authy-openvpn)
which isn't updated anymore and is lacking Authy OneTouch support.

## Pre-Requisites

1.  Authy API key
2.  OpenVPN server installation

## Installation

### From compiled package

1.  Download latest [release](https://github.com/3fs/go-authy-openvpn/releases) from GitHub
2.  Extract tar archive in desired location on your server
3.  Run `post-install` script

### From source

1.  Install requirements: Golang, build essentials
2.  Run `sudo make install`

### Check OpenVPN config

Open your OpenVPN configuration file in your prefered editor.
If you allowed `post-instal` script from previous step to edit your `server.conf`
you should now have something like this at the end of it:

```
plugin <go-authy-path>/auth_script.so <go-authy-path>/go-authy-openvpn -a <authy_api_key>
```

If you don't have something like tihs in your configuration, you should add it manually.

If you want to have authy configuration file in non-defautl location (`/etc/openvpn/authy/authy-vpn.conf`)
you can add `-c /path/to/authy-vpn.conf` after your api key:

```
plugin <go-authy-path>/auth_script.so <go-authy-path>/go-authy-openvpn -a <authy_api_key> -c /path/to/authy-vpn.conf
```

## Migrating from official plugin

If you are already using official [Authy OpenVPN plugin](https://github.com/authy/authy-openvpn) you can
install this plugin with the steps above and then remove the old plugin from your OpenVPN server config.
This plugin uses the same format of `authy-vpn.conf` so all your registered users will stay there.

## Adding users

This plugin comes with a script, that helps you register users.
To start adding users type: `sudo authy-vpn-add-user`

If the script was successful it will add username and Authy ID to `/etc/openvpn/authy/authy-vpn.conf`.

## How it works

This plugin works with certificates based authentication. To login the user need certificate, username and password.

Password can be 4 different things:

### OneTouch

If the provided password is `onetouch` the user will receive OneTouch push notification to Authy app where they can approve the login.

### Token

Password can be TOTP token from Authy app or token the user received through SMS or call.

### SMS or call

If the password is `sms` or `call` the plugin will make a request for that to Authy and will fail the login.
Then the user will receive the token through SMS or call and will then use that token on next login.

## VPN client configuration

Your users will need to add

```
auth-user-pass
```

to their `client.conf`.
