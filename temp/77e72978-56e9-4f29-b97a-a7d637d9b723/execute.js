//import ./code.js and write its output to file
var ret = require('./code.js');
var fs = requir('fs');

fs.writeFileSync('./out.txt', ret);
