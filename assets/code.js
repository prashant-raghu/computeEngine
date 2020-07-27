function code(){
	for(let i=0;i<10000000000;i++){}
	return "Hello, this originates from code.js"
}
module.exports = code();