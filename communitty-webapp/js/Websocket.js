define(['underscore'], function(_) {
    var Websocket = function(url) {
        this.url = url;

        // Initialize some sets of handlers for
        // the various events
        this.onErrorHandlers = [];
        this.onOpenHandlers = [];
        this.onCloseHandlers = [];
        this.onMessageHandlers = [];

        this.connection = new WebSocket(url);
        this.connection.onerror = _.bind(this._onError, this);
        this.connection.onopen = _.bind(this._onOpen, this);
        this.connection.onclose = _.bind(this._onClose, this);
        this.connection.onmessage = _.bind(this._onMessage, this);

    };

    Websocket.prototype = {
        /**
         *
         */
        onError: function(handler) {
            this.onErrorHandlers.push(handler);
        },

        /**
         *
         */
        onOpen: function(handler) {
            this.onOpenHandlers.push(handler);
        },

        /**
         *
         */
        onClose: function(handler) {
            this.onCloseHandlers.push(handler);
        },

        /**
         *
         */
        onMessage: function(handler) {
            this.onMessageHandlers.push(handler);
        },

        /**
         *
         */
        _onError: function(event) {
            for (var i in this.onErrorHandlers) {
                this.onErrorHandlers[i](event);
            }
        },

        /**
         *
         */
        _onOpen: function(event) {
            for (var i in this.onOpenHandlers) {
                this.onOpenHandlers[i](event);
            }
        },

        /**
         *
         */
        _onClose: function(event) {
            for (var i in this.onCloseHandlers) {
                this.onCloseHandlers[i](event);
            }
        },

        /**
         *
         */
        _onMessage: function(event) {
            for (var i in this.onMessageHandlers) {
                this.onMessageHandlers[i](event);
            }
        }
    };

    return Websocket;
});
