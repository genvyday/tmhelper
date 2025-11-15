//while(!tmh.CptKey(tmh.Pwd("crypto key:"),"encrypted_123_for_check"))
//  tmh.Println("crypto key not match");
tmh.CptKey("12345","GY2dDlxuHET_U0NqwpZyzQ");
const creds=new Map([ //credentials
  ["devdrt",{user:"devuser",pwd:"ejXJ3n82nxpKp6IVYd-f4Q"}],
  ["uatbst",{user:"uatbst",pwd:"ejXJ3n82nxpKp6IVYd-f4Q"}],
  ["prdbst",{user:tmh.Dec("jBcdXgsKVZHEIOxvV0LwNQ"),pwd:"QfD4JR0V0ySRWjL4sIKCNQ"}],
  ["uatweb",{user:"webuser",pwd:"QfD4JR0V0ySRWjL4sIKCNQ"}],
  ["prdweb",{user:"prdweb",pwd:"QfD4JR0V0ySRWjL4sIKCNQ"}],
  ["prdsvc",{user:"prdsvc",pwd:"QfD4JR0V0ySRWjL4sIKCNQ"}],
]);

const hosts=[
    {name:"exit   ",host:"for exit"}
	{name:"dev.001",host:"172.16.0.2",cred:creds.get("devdrt"),lgn:directlgn},
	{name:"uat.chn",host:"192.168.0.2",cred:creds.get("uatbst"),lgn:bstlgn,scred:uatscred},
	{name:"prd.usa",host:tmh.Dec("bnP6wmuI9bi3l9dnn1Y6cg"),cred:creds.get("prdbst"),lgn:bstlgn,scred:prdscred,port:"22"},
];
const hcount=hosts.length;
const prmt="Please Select[0-"+hcount+"]：";
tmh.SetTimeout(0);
var tp=tmh.NewTerm();
mainloop();
function mainloop()
{
	var hidx=hcount;
	while(hidx!=0)
	{
		if(tmh.Goos()=="windows") tp.WaitDone("\nPress Any Key To Select Host:");
		tmh.Println("\n\nHost List：");
		for(i=0;i<hcount;++i)
		{
			tmh.Println("\t"+i+": "+hosts[i].name+"\t"+hosts[i].host);
		}
		hidx=hcount;
		for(j=0;j<10&&(hidx>=hcount||isNaN(hidx));++j)
		{
			s=tp.Input(prmt);
			hidx=parseInt(s);
			if(hidx<hosts.length&&hidx!=0) tp=hosts[hidx].lgn(hosts[hidx]);
		}
	}
}
function bstlgn(host) //login bastion host first
{
    tp=sshconnect(host);
	tp.Matchs([[") Password:", tmh.Dec(host.cred.pwd)+"\n"]]);
	while(tp.Ok())
	{
		tp.Expect("taget host:"); //login application host
		s=tp.ReadStr("\n");       //read user input from server echo
		cred=host.scred(parseInt(s)); //get credential
		if(cred!=null)
		{
			tp.Matchs([["login:",cred.user+"\n","C"],["password:",tmh.Dec(cred.pwd)+"\n"]]);
		}
		else
		{
			tmh.Println(s);
		}
	}
	return tp.Exit();
}
function directlgn(host) //login application host direct
{
    tp=sshconnect(host);
	cred=host.cred
	tp.Matchs([["login:",cred.user+"\n","C"],["password:",tmh.Dec(cred.pwd)+"\n"]]);
	return tp.Term();
}
function sshconnect(host)
{
    port="22";
    if(host.port!=null) port=host.port;
	return tp.Exec(["ssh",host.cred.user+"@"+host.host,"-p",port]);
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
