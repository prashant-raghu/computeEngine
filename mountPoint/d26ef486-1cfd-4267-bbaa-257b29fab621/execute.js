//import ./code.js and write its output to file
var ret = require('./code.js');
var fs = require('fs');
try {
    console.log('No err so far')
    fs.writeFileSync('/app/out.txt', ret);
}
catch (err) {
    console.log(err)
}
