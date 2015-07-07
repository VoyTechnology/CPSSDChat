var config = require('./config.json');
var App = require("./System/app.js");

var app = new App();

app.start( config );
