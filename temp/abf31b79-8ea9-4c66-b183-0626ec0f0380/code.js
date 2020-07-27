function code(){
    let i = 0;
    while(i++){if(i==100000000000000000)break;}
	 return "Hello, this originates from code.js after billion iterations"
}
module.exports = code();