CommuniTTY shares your TTY via the browser so people can see what you're up to.

Building
========

To build the service, run

    make deps
    make

You only need to run `make deps` once.

Running
=======

To run the service, build and then run

    ./communitty

It'll respawn your shell and serve your TTY via

    http://localhost:9000/

The served data is read-only.
