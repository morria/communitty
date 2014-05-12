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

require(['jquery', 'underscore', 'terminal', 'TerminalSocket'],
  function($, _, Terminal, TerminalSocket) {
    new TerminalSocket(
      'ws://' + location.host + '/term',
      $('#terminal'),
      $('#messages'));
  }
);
