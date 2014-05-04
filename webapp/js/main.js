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

    // connection.binaryType = 'arraybuffer';
    // connection.binaryType = 'blob';

    connection.onerror = function(event) {
        console.err(event);
    }

    connection.onopen = function(event) {
        console.log('open');
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

    console.log(connection);

    connection.onmessage = _.bind(function(message) {
      console.log(message);
      terminal.write(message.data);
    }, this);

  }
);
