import{d as S,k as p,w,U as h,V as C,f as l,$ as s,a8 as D,a4 as i,_,u as e,E as m,a3 as N}from"./vue-14860272.js";import{af as f}from"./index-6940427b.js";import"./dayjs-4778c158.js";import{u as g}from"./useModal-231cd172.js";import{D as k,p as V}from"./index-00d3628b.js";import{V as I}from"./jb4-6f943bfa.js";import{i as P,a as y}from"./gb28281-bee701fe.js";import{a as B,D as u}from"./index-8d59d51b.js";import{C as z}from"./Card-8c66b56c.js";const $={class:"h-full w-full flex flex-items-center"},E={class:"text-center"},T=S({__name:"play",props:{path:null,id:null,deviceName:null,deviceID:null,channelName:null},setup(a){const n=a,d={fontWeight:"bold"},b=p("http-flv"),o=p("");w([()=>n.id,()=>n.path],([t,c])=>{c?o.value=c:t&&P({id:n.id,channel:n.deviceID}).finally(()=>{o.value=`${n.id}/${n.deviceID}`})},{immediate:!0});const r=()=>{const t=o.value.split("/");return{id:t[0],channel:t[1]}},v=async t=>{f.destroy(),await y({id:r().id,channel:r().channel,ptzcmd:t}).then(async c=>{f.success("指令发送成功"),x()}).catch(c=>{f.success("指令发送失败")})},x=()=>{setTimeout(async()=>{await y({id:r().id,channel:r().channel,ptzcmd:V()})})};return(t,c)=>(h(),C("div",$,[l(I,{videoShadow:!0,streamPath:o.value,format:b.value,class:"flex-1 m-r-10px"},{default:s(()=>[D("div",E,i(o.value),1)]),_:1},8,["streamPath","format"]),l(e(z),{bordered:!1,class:"dark:bg-transparent",bodyStyle:{height:"100%",width:"280px"}},{default:s(()=>[l(k,{hanldeClick:v,allowed:!!o.value},null,8,["allowed"]),a.id?(h(),_(e(B),{key:0,column:1,class:"mt-20px",layout:"vertical"},{default:s(()=>[l(e(u),{class:"!pb-8px",labelStyle:e(d),label:"设备名称"},{default:s(()=>[m(i(a.deviceName||"--"),1)]),_:1},8,["labelStyle"]),l(e(u),{class:"!pb-8px",labelStyle:e(d),label:"设备编号"},{default:s(()=>[m(i(a.id||"--"),1)]),_:1},8,["labelStyle"]),l(e(u),{class:"!pb-8px",labelStyle:e(d),label:"通道名称"},{default:s(()=>[m(i(a.channelName||"--"),1)]),_:1},8,["labelStyle"]),l(e(u),{class:"!pb-8px",labelStyle:e(d),label:"通道编号"},{default:s(()=>[m(i(a.deviceID||"--"),1)]),_:1},8,["labelStyle"])]),_:1})):N("",!0)]),_:1})]))}}),H=a=>{g({content:()=>l(T,a,null),modalConfig:{width:"100%",wrapClassName:"full-antdv-modal",destroyOnClose:!0,footer:null}})};export{H as u};