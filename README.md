# Requirement: #

Squid 3.4

# Get API Key: #

- Login to your Google account
- Point your browser to https://developers.google.com/safe-browsing/key_signup
- Accept the terms and condition and generate your API key

# Install: #

Just download the appropiate file as simple as that.

    $ sudo curl -o /usr/local/bin/squid-gsb https://raw.githubusercontent.com/catinello/squid-google-safe-browsing/master/bin/squid-gsb-$(uname -m)

Supported pre-compiled architectures are PC-64bit (x86_64) and Raspberry Pi (armv61). I'm just using and testing on those 2 architectures. It is compilable for others though.

# Usage: #

Get the amount of cpu cores available:

    $ grep ^processor /proc/cpuinfo | tail -n 1 | awk -F': ' '{print $2}'

Use this number eg. 4 (on a dual-core with hyperthreading) for the following concurrency setting.

edit /etc/squid/squid.conf

    url_rewrite_children 20 startup=0 idle=1 concurrency=4
    url_rewrite_program /usr/local/bin/squid-gsb [GSB_APIKEY]

# Environment Variable: #

It is possible to use an environment variable to store the API key.

    $ export GSB_APIKEY=WHATEVERYOURAPIKEYIS

# Logging: #

All errors and blocks are being logged to your syslog facility.

A blocked and redirected access looks like this:

    Apr 08 17:35:32 03-proxy squid-gsb[28473]: Blocked Site: http://ianfette.org

A hint that your API key is invalid:

    Apr 07 15:47:48 03-proxy squid-gsb[27416]: Not Authorized


# Background: #

I know that this feature is already build-in in chrome and firefox. Still there are other browsers out there and my main reason was to protect people from themselfs.

A user can't work around a site warning from chrome or firefox. If there is really a false positive you can always whitelist it on the squid proxy.

I have no affiliation with Google.

# License: #

[&copy; Antonino Catinello][HOME] - [MIT-License][MIT]

[MIT]:https://github.com/catinello/squid-google-safe-browsing/blob/master/LICENSE
[HOME]:http://antonino.catinello.eu
