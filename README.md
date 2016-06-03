# Deprecation Notice: #

Google will stop supporting the v3 API in early 2017 and I abandoned this project long ago.

# Requirement: #

- Squid 3.4
- Safe-Browsing-API v3 Key

# Get API Key: #

- Login to your Google account
- Point your browser to https://developers.google.com/safe-browsing/key_signup
- Accept the terms and condition and generate your API key

# Compile: #

    $ cd $GOPATH
    $ go get -d github.com/catinello/squid-google-safe-browsing
    $ cd src/github.com/catinello/squid-google-safe-browsing
    $ go build -o squid-gsb

# Usage: #

Get the amount of cpu cores available:

    $ grep ^processor /proc/cpuinfo | tail -n 1 | awk -F': ' '{print $2}'

Use this number eg. 4 (on a dual-core with hyperthreading) for the following concurrency setting.

edit /etc/squid/squid.conf

    url_rewrite_children 20 startup=0 idle=1 concurrency=4
    url_rewrite_program /usr/local/bin/squid-gsb [GSB_APIKEY]

# Environment Variable: #

It is possible to use an environment variable to store the API key for testing purposes.

    $ export GSB_APIKEY=WHATEVERYOURAPIKEYIS

# Debugging: #

Create the following file to enable debugging which is checked pre loop which means you eventually have to restart the service.

    $ touch /tmp/squid-gsb.debug
    $ sudo systemctl restart squid #optional

You would get the following output in your syslog (info): url -> gsb-result: channel-id squid-result-code

    Aug 23 17:19:59 03-proxy squid-gsb[28523]: http://www.google.com:443 -> 204: 0 ERR

# Logging: #

All errors and blocks are being logged to your syslog facility (critical).

A blocked and redirected access looks like this:

    Apr 08 17:35:32 03-proxy squid-gsb[28473]: Blocked Site: http://ianfette.org

A hint that your API key is invalid:

    Apr 07 15:47:48 03-proxy squid-gsb[27416]: Not Authorized

# Background: #

I know that this feature is already build-in in chrome and firefox. Still there are other browsers out there and my main reason is to protect people from themselfs.

This way a user can't work around a site warning as in chrome or firefox. If there is a false positive, then you can simple whitelist it through the squid proxy configuration.

I have no affiliation with Google.

# License: #

[&copy; Antonino Catinello][HOME] - [MIT-License][MIT]

[MIT]:https://github.com/catinello/squid-google-safe-browsing/blob/master/LICENSE
[HOME]:http://antonino.catinello.eu
