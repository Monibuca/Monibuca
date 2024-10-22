import{d as _,r as B,U as f,V as v,f as o,W as h,a8 as k,u as m,$ as z,E as I,Y as A,o as L,w as C,k as x,F as U,a7 as T,_ as P,al as R}from"./vue-14860272.js";import{b as M,c as E}from"./global-3f00abc7.js";import"./dayjs-4778c158.js";import{aj as K,q as V,aC as W}from"./index-6940427b.js";import{S as H}from"./index-e0c9a531.js";import{_ as O}from"./index.vue_vue_type_script_setup_true_lang-8aad6e65.js";import{_ as F,S,l as G,d as Z,a as j,b as Y}from"./hard-disk.vue_vue_type_script_setup_true_lang-63f30946.js";import{g,a as N,b,c as y}from"./utils-bc030ba0.js";import"./mockjs-890b569b.js";const q={class:"flex"},J=_({__name:"top-info",props:{data:null,hardDisk:null},setup(u){const t=B({Version:"",StartTime:"",LocalIP:""});M().then(r=>{t.LocalIP=r.LocalIP,t.Version=r.Version,t.StartTime=K.toDateString(new Date(r.StartTime),"yyyy-MM-dd HH:ss:mm")}).catch(r=>{console.error(`getSysInfo-err: ${r}`)});const e="p-14px rounded-16px bg-#fff dark:bg-#100C2A flex-center flex-1 mr-10px";return(r,n)=>{var a;return f(),v("div",q,[o(F,{class:h([e,"min-w-320px"]),hardDisk:u.hardDisk},null,8,["hardDisk"]),k("div",{class:h([e,"min-w-140px"])},[o(m(H),{icon:"streams",size:"70",class:"text-#bb86fc"}),o(m(S),{title:"当前流数",value:(a=u.data)==null?void 0:a.length},null,8,["value"])]),k("div",{class:h([e,"min-w-210px flex-col flex-items-start"])},[o(m(S),{title:"本地IP",value:t.LocalIP,valueStyle:{fontSize:"18px",textAlign:"left"}},null,8,["value"]),o(m(S),{title:"启动时间",value:t.StartTime,valueStyle:{fontSize:"18px"}},{suffix:z(()=>[I(" ["),o(m(O),{value:t.StartTime,class:"text-primary"},null,8,["value"]),I("] ")]),_:1},8,["value"])]),k("div",{class:h([e,"min-w-180px !mr-0px"])},[o(m(S),{class:"overflow-auto",title:"当前版本",value:t.Version,valueStyle:{fontSize:"18px"}},null,8,["value"])])])}}}),Q=["id"],X=_({__name:"network-item",props:{network:null},setup(u){const t=u,e={text:t.network.Name||"",receiveData:[g(t.network.ReceiveSpeed)],sentData:[g(t.network.SentSpeed)],timeData:[N()],sent:t.network.Sent,receive:t.network.Receive},r={backgroundColor:"",title:{text:e.text,top:-5},tooltip:{trigger:"axis",formatter:function(s){var c;var l=s[0].name+"<br>";for(let p of s)l+=((c=p.seriesName)==null?void 0:c.split(":")[0])+" : "+p.value+" KB/s <br>";return l}},legend:{data:[`发送: ${y(e.sent)}`,`接收: ${y(e.receive)}`],right:10},axisPointer:{link:{xAxisIndex:"all"}},grid:[{left:50,right:50,height:"35%"},{left:50,right:50,top:"55%",height:"35%"}],xAxis:[{type:"category",boundaryGap:!1,axisLine:{onZero:!0},data:e.timeData},{gridIndex:1,type:"category",boundaryGap:!1,axisLine:{onZero:!0},data:e.timeData,position:"top"}],yAxis:[{name:"发送 (KB/s)",type:"value"},{name:"接收 (KB/s)",gridIndex:1,type:"value",inverse:!0}],series:[{name:`发送: ${y(e.sent)}`,type:"line",smooth:!0,showSymbol:!0,data:e.sentData,label:{show:!0},lineStyle:{color:"#bb86fc"},itemStyle:{color:"#bb86fc"}},{name:`接收: ${y(e.receive)}`,type:"line",xAxisIndex:1,yAxisIndex:1,showSymbol:!0,data:e.receiveData,label:{show:!0},lineStyle:{color:"#5a00ff"},itemStyle:{color:"#5a00ff"}}]},{getDarkMode:n}=A(V());let a,i;L(()=>{d(!0)});const d=(s=!1)=>{s||(i==null||i.stop(),a&&Z(a)),a=G(t.network.Name,r,m(n)),i=W(document.body,()=>{a&&(a==null||a.resize())})};return C(()=>t.network,s=>{e.text=s.Name,e.timeData=b(e.timeData,N()),e.sentData=b(e.sentData,g(s.SentSpeed)),e.receiveData=b(e.receiveData,g(s.ReceiveSpeed)),e.sent=s.Sent,e.receive=s.Receive,a.setOption(r)},{immediate:!1}),C(()=>m(n),()=>d(!1)),(s,l)=>(f(),v("div",{id:u.network.Name,class:"rounded-16px bg-#fff dark:bg-#100C2A h-328px p-14px"},null,8,Q))}}),ee={class:"flex flex-wrap mt-8px"},te=_({__name:"network",setup(u,{expose:t}){const e=x([]);return t({handleUpdate:n=>{e.value=n}}),(n,a)=>(f(),v("div",ee,[(f(!0),v(U,null,T(e.value,i=>(f(),P(X,{key:i.Name,network:i,class:"flex-1 min-w-380px mr-8px mb-8px"},null,8,["network"]))),128))]))}}),ae={class:"mt-8px flex"},se=_({name:"Overview"}),fe=_({...se,setup(u){const t=x(),e=x(),r=x(),n=x();let a;return(()=>{a=E(d=>{var s,l,c,p,w,D,$;t.value=d,(c=e==null?void 0:e.value)==null||c.handleUpdate((l=(s=t.value)==null?void 0:s.Memory)==null?void 0:l.Usage),(w=r==null?void 0:r.value)==null||w.handleUpdate((p=t.value)==null?void 0:p.CPUUsage),($=n==null?void 0:n.value)==null||$.handleUpdate((D=t.value)==null?void 0:D.NetWork)})})(),R(d=>{const{path:s}=d;s!=="/stream-push/list"&&(a==null||a())}),(d,s)=>{var l,c;return f(),v(U,null,[o(J,{hardDisk:(l=t.value)==null?void 0:l.HardDisk,data:((c=t.value)==null?void 0:c.Streams)||[]},null,8,["hardDisk","data"]),k("div",ae,[o(j,{ref_key:"memoryRef",ref:e,class:"flex-1 mr-10px"},null,512),o(Y,{ref_key:"cpuRef",ref:r,class:"flex-1"},null,512)]),o(te,{ref_key:"networkRef",ref:n},null,512)],64)}}});export{fe as default};