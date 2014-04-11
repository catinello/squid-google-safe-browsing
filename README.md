# Requirement: #

Squid 3.4

# Get API Key: #

- Login to your Google account
- Point your browser to https://developers.google.com/safe-browsing/key_signup
- Accept the terms and condition and generate your API key

# Usage: #

edit /etc/squid/squid.conf

    url_rewrite_program /usr/local/bin/squid-gsb-x86_64 GSB_APIKEY

# Environment Variable: #

It is possible to use an environment variable to store the API key.

    $ export GSB_APIKEY=WHATEVERYOURAPIKEYIS

# Logging: #

All errors and blocks are being logged to your syslog facility.

A blocked and redirected access looks like this:

    Apr 08 17:35:32 03-proxy squid-gsb[28473]: Blocked Site: http://ianfette.org

A hint that your API key is invalid:

    Apr 07 15:47:48 03-proxy squid-gsb[27416]: Not Authorized


# ToDo: #

- Use concurrency from squid by default.

# Background: #

I know that this feature is already build-in in chrome and firefox. Still there are other browsers out there and my main reason was to protect people from themselfs.

A user can't work around a site warning from chrome or firefox. If there is really a false positive you can always whitelist it on the squid proxy.

I have no affiliation with Google.

# License: #

[&copy; Antonino Catinello][HOME] - [MIT-License][MIT]

[MIT]:https://github.com/catinello/squid-google-safe-browsing/blob/master/LICENSE
[HOME]:http://antonino.catinello.eu
