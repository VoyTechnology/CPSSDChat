var http = require('http');

/**
 * Main Chat Application
 * @class
 **/
var App = {

  start: function( config ){
    http.createServer( function( req, res ){

    }).listen(config.server.port);
  }

};
