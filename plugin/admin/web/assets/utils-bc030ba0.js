import{aj as a}from"./index-6940427b.js";import"./vue-14860272.js";const u=()=>a.toDateString(new Date,"HH:ss:mm");function o(t,r){return t.length>=20&&t.shift(),t.push(r),t}const f=t=>Math.round(t/1024*100)/100,h=t=>t<1024?`${t} B`:(t=t/1024,t<1024?`${Math.round(t*100)/100} KB`:(t=t/1024,t<1024?`${Math.round(t*100)/100} MB`:`${Math.round(t*100)/100} GB`));export{u as a,o as b,h as c,f as g};
