define(['underscore', 'jquery', 'Websocket'],
function (_, $, Websocket) {
    var TerminalSocket = function(url, $terminal, $messages) {
        // The URL of the terminal endpoint
        this.url = url;
        this.$terminal = $terminal;
        this.$messages = $messages;

        this.terminal = new Terminal({
            cols: 80,
            rows: 24,
            screenKeys: false
        });
        this.terminal.on('title', _.bind(this.onTerminalTitle, this));
        this.terminal.open($terminal.get(0));

        this.socket = null;
        this.reconnectTimer = null;

        this.connect();
    };

    TerminalSocket.prototype = {
        connect: function() {
            if (this.socket && this.socket.hasOwnProperty('close')) {
                this.socket.close();
                this.socket = null;
            }

            this.socket = new Websocket(this.url);
            this.socket.onError(_.bind(this.onSocketError, this));
            this.socket.onMessage(_.bind(this.onSocketMessage, this));
            this.socket.onClose(_.bind(this.onSocketClose, this));
            this.socket.onOpen(_.bind(this.onSocketOpen, this));
        },

        /**
         *
         */
        onSocketError: function(event) {
            // console.error("socket error", event)
        },

        /**
         * When we lose our connection, attempt to reestablish
         * one in awhile.
         */
        onSocketClose: function(event) {
            // Notify the user that we're disconnected
            this.$messages.find('#disconnected').show()

            // n.b.: Websocket calls onclose when a
            //       connection attempt fails, so we
            //       only need to set a single timer.
            this.reconnectTimer = setTimeout(_.bind(function(event) {
                this.connect();
            }, this), 500);
        },

        /**
         *
         */
        onSocketOpen: function(event) {
            // Remove any disconnected notices
            this.$messages.find('.message').hide()
            this.$terminal.show();

            // Kill any reconnect timers that may be
            // kicking around
            if (this.reconnectTimer) {
                clearInterval(this.reconnectTimer);
                this.reconnectTimer = null;
            }
        },

        /**
         *
         */
        onSocketMessage: function(event) {
            var command = JSON.parse(JSON.parse(event.data))

            switch (command.Command) {
                case "Terminal":
                    this.terminal.write(atob(command.Data));
                    break;
                case "WindowSize":
                    this.terminalResize(command.Cols, command.Rows);
                    break;
            }
        },

        terminalResize: function(cols, rows) {
            this.terminal.resize(cols, rows);
        },

        /**
         *
         */
        onTerminalTitle: function(title) {
            document.title = title;
        },
    };

    return TerminalSocket;
});
