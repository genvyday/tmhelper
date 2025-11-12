//tmh_crypto(tmh_pwd("crypto key:"));
tmh_crypto("12345")
const creds=new Map([ //credentials
  ["devdrt",{user:"devuser",pwd:"ejXJ3n82nxpKp6IVYd-f4Q"}],
  ["uatbst",{user:"uatbst",pwd:"ejXJ3n82nxpKp6IVYd-f4Q"}],
  ["prdbst",{user:tmh_dec("jBcdXgsKVZHEIOxvV0LwNQ"),pwd:"QfD4JR0V0ySRWjL4sIKCNQ"}],
  ["uatweb",{user:"webuser",pwd:"QfD4JR0V0ySRWjL4sIKCNQ"}],
  ["prdweb",{user:"prdweb",pwd:"QfD4JR0V0ySRWjL4sIKCNQ"}],
  ["prdsvc",{user:"prdsvc",pwd:"QfD4JR0V0ySRWjL4sIKCNQ"}],
]);

var hosts=[
	{name:"dev.001",host:"172.16.0.2",cred:creds.get("devdrt"),lgn:directlgn},
	{name:"uat.chn",host:"192.168.0.2",cred:creds.get("uatbst"),lgn:bstlgn,scred:uatscred},
	{name:"prd.usa",host:tmh_dec("bnP6wmuI9bi3l9dnn1Y6cg"),cred:creds.get("prdbst"),lgn:bstlgn,scred:prdscred,port:"22"},
];
tmh_setTimeout(999999999);
tmh_println("\n\nhosts list：");
for(i=0;i<hosts.length;++i)
{
    tmh_println(i+": "+hosts[i].name+"\t"+hosts[i].host);
}
i=tmh_input("select host:");
var host=hosts[parseInt(i)];
if(host!=null) host.lgn(host);
function bstlgn(host) //login bastion host first
{
    sshconnect(host);
	tmh_matchs([[") Password:", tmh_dec(host.cred.pwd)+"\n"]]);
	while(tmh_ok())
	{
		tmh_expect("taget host:"); //login application host
		s=tmh_readStr("\n");       //read user input from server echo
		cred=host.scred(parseInt(s)); //get credential
		if(cred!=null)
		{
			tmh_matchs([["login:",cred.user+"\n","C"],["password:",tmh_dec(cred.pwd)+"\n"]]);
		}
		else
		{
			tmh_println(s);
		}
	}
	tmh_exit();
}
function directlgn(host) //login application host direct
{
    sshconnect(host);
	cred=host.cred
	tmh_matchs([["login:",cred.user+"\n","C"],["password:",tmh_dec(cred.pwd)+"\n"]]);
	tmh_term();
}
function sshconnect(host)
{
    port="22";
    if(host.port!=null) port=host.port;
	tmh_run(["ssh",host.cred.user+"@"+host.host,"-p",port]);
}
function uatscred(n) //map to credential
{
	if(n>3) return creds.get("uatweb");
	else return null;
}
function prdscred(n) //map to credential
{
	if(n<7) return creds.get("prdweb");
	else if(n>10&&n<16) return creds.get("prdsvc");
	else return null;
}
