require.config({
  shim: {
    underscore: {
      exports: '_'
    },
    terminal: {
      exports: 'Terminal'
    }
  },
  paths: {
    underscore: 'vendor/underscore-min',
    jquery: 'vendor/jquery-2.0.3.min',
    terminal: 'vendor/term'
  }
});

require(['jquery', 'underscore', 'terminal'],
  function($, _, Terminal) {
    var connection =
        new WebSocket('ws://' + location.host + '/term');

    connection.onerror = function(event) {
        console.log("onerror: ", event);
    }

    connection.onopen = function(event) {
    }

    var terminal = new Terminal({
        cols: 80,
        rows: 24,
        screenKeys: false
    });

    terminal.on('title', function(title) {
        document.title = title;
    });

    terminal.open($('#terminal').get(0));

    connection.onmessage = _.bind(function(message) {
      // Yeah, we're doing this
      var command = JSON.parse(JSON.parse(message.data))

      if ("WindowSize" === command.Command) {
        terminal.resize(command.Cols, command.Rows);

      } else if ("Terminal" === command.Command) {
        terminal.write(atob(command.Data));
      }

    }, this);

  }
);
