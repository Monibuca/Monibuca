import{d as h,f as e,e as T,k as v,w as V,U as y,_ as x,$ as s,a8 as k,u as t,E as g,V as P,F as O,a7 as E,a4 as L,a3 as F}from"./vue-14860272.js";import{g as R,a as $,u as z}from"./global-3f00abc7.js";import{u as q}from"./formily-de1f5b97.js";import"./dayjs-4778c158.js";import{u as G,d as S,a as w,P as M,M as J,X as N,I as j,a9 as A,J as B,ae as U,af as X}from"./index-6940427b.js";import{C as p}from"./Card-8c66b56c.js";import{A as H}from"./index-4555d001.js";import"./index-8bf1b192.js";import"./LeftOutlined-5bd8a798.js";import"./isNumeric-3f69e2aa.js";import"./index-9788fdab.js";import"./index-5174f10d.js";import"./index-00d65388.js";import"./index-14bb9a5e.js";import"./Group-70d09d1e.js";import"./useFlexGapSupport-d0720122.js";import"./index-2c49b025.js";import"./index-3d8f6fc8.js";import"./index-5ffcd845.js";import"./scrollTo-c1f6ee63.js";import"./mockjs-890b569b.js";var Q=function(){return{prefixCls:String,title:M.any,description:M.any,avatar:M.any}};const D=h({compatConfig:{MODE:3},name:"ACardMeta",props:Q(),slots:["title","description","avatar"],setup:function(l,f){var n=f.slots,m=G("card",l),r=m.prefixCls;return function(){var c=S({},"".concat(r.value,"-meta"),!0),o=w(n,l,"avatar"),_=w(n,l,"title"),b=w(n,l,"description"),i=o?e("div",{class:"".concat(r.value,"-meta-avatar")},[o]):null,d=_?e("div",{class:"".concat(r.value,"-meta-title")},[_]):null,a=b?e("div",{class:"".concat(r.value,"-meta-description")},[b]):null,C=d||a?e("div",{class:"".concat(r.value,"-meta-detail")},[d,a]):null;return e("div",{class:c},[i,C])}}});var W=function(){return{prefixCls:String,hoverable:{type:Boolean,default:!0}}};const I=h({compatConfig:{MODE:3},name:"ACardGrid",__ANT_CARD_GRID:!0,props:W(),setup:function(l,f){var n=f.slots,m=G("card",l),r=m.prefixCls,c=T(function(){var o;return o={},S(o,"".concat(r.value,"-grid"),!0),S(o,"".concat(r.value,"-grid-hoverable"),l.hoverable),o});return function(){var o;return e("div",{class:c.value},[(o=n.default)===null||o===void 0?void 0:o.call(n)])}}});p.Meta=D;p.Grid=I;p.install=function(u){return u.component(p.name,p),u.component(D.name,D),u.component(I.name,I),u};const Y={class:"flex h-full"},Z={class:"w-700px"},ee={class:"flex justify-between p-4px",style:{background:"var(--layout-background)"}},ae={key:1,class:"absolute top-40% left-50%"},te=h({name:"Config"}),we=h({...te,setup(u){const l=v([]);R().then(i=>{Object.keys(i).length&&(l.value=Object.entries(i).map(([d,a])=>({name:a.Name,version:`版本: ${a.Version}`,disabled:a.Disabled})))});const f=v(null),n=v(),m=v(!1),r=v(["global-sub","plugin-sub"]),c=v(["global"]);V(c,i=>{const d={name:i[0]==="global"?"":i[0],formily:"1"};m.value=!0,$(d).then(a=>{const{FormilyForm:C,form:K}=q({schema:a});f.value=C,n.value=K}).finally(()=>{m.value=!1})},{immediate:!0});const o=async()=>{if(!n.value.modified)return X.warning("您尚未修改过任何配置，无需保存！");const i=await n.value.submit();z({name:c.value[0]},i)},_=()=>document.querySelector(".plugin-list"),b=()=>{window.open("https://monibuca.com/docs/guide/plugins/"+c.value[0].toLowerCase()+".html","_blank")};return(i,d)=>(y(),x(t(p),{class:"flex-center h-100% overflow-hidden",bodyStyle:{height:"100%"}},{default:s(()=>[k("div",Y,[e(t(J),{selectedKeys:c.value,openKeys:r.value,mode:"inline",inlineCollapsed:!1,onSelect:d[0]||(d[0]=({key:a})=>c.value=[a]),class:"w-180px h-100% overflow-y-auto overflow-x-hidden"},{default:s(()=>[e(t(N),{key:"global-sub"},{title:s(()=>[g("全局配置")]),default:s(()=>[e(t(j),{key:"global"},{default:s(()=>[g(" global ")]),_:1})]),_:1}),e(t(N),{key:"plugin-sub"},{title:s(()=>[g("插件配置")]),default:s(()=>[(y(!0),P(O,null,E(l.value,a=>(y(),x(t(j),{key:a.name,disabled:a.disabled},{default:s(()=>[g(L(a.name),1)]),_:2},1032,["disabled"]))),128))]),_:1})]),_:1},8,["selectedKeys","openKeys"]),e(t(p),{bordered:!1,bodyStyle:{paddingTop:0},class:"plugin-list h-full overflow-y-auto overflow-x-hidden relative"},{default:s(()=>[k("div",Z,[e(t(H),{offsetTop:0,target:_,class:"mb-10px"},{default:s(()=>[k("div",ee,[e(t(A),{type:"link",onClick:o},{default:s(()=>[e(t(B),{icon:"tabler:hand-click",class:"v-text-bottom"}),g(" 点我保存配置 ")]),_:1}),e(t(A),{type:"link",onClick:b},{default:s(()=>[g(" 配置文档 "),e(t(B),{icon:"fluent:window-new-16-filled",class:"v-text-bottom"})]),_:1})])]),_:1}),f.value?(y(),x(t(f),{key:0})):F("",!0),m.value?(y(),P("div",ae,[e(t(U),{class:"zIndex-100",size:"large"})])):F("",!0)])]),_:1})])]),_:1}))}});export{we as default};