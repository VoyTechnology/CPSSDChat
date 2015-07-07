var http = require('http');

/**
 * Main Chat Application
 * @class
 **/
var App = function(){}

App.prototype = {
  start: function(config){
    http.createServer( function( req, res ){
      res.writeHead(200, {'Content-Type':'application/json'});
      res.end('{"hello":"world"}');
    }).listen(config.server.port);

    console.log("Server Started on port "+ config.server.port);
  }
}

module.exports = App;
