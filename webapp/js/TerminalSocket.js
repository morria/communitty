define(['underscore', 'jquery', 'Websocket'],
function (_, $, Websocket) {
    var TerminalSocket = function(url, $terminal) {
        this.socket = new Websocket(url);
        this.socket.onError(_.bind(this.onSocketError, this));
        this.socket.onMessage(_.bind(this.onSocketMessage, this));

        this.terminal = new Terminal({
            cols: 80,
            rows: 24,
            screenKeys: false
        });
        this.terminal.on('title', _.bind(this.onTerminalTitle, this));
        this.terminal.open($terminal.get(0));
    };

    TerminalSocket.prototype = {
        /**
         *
         */
        onSocketError: function(event) {
            console.error(event)
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
            console.log("resize", rows, cols);
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
