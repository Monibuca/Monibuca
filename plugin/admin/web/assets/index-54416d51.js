import{l as p,u as S,d as a,b as d,f as D}from"./index-6940427b.js";import{d as _,e as r,f as v}from"./vue-14860272.js";var b=function(){return{prefixCls:String,type:{type:String,default:"horizontal"},dashed:{type:Boolean,default:!1},orientation:{type:String,default:"center"},plain:{type:Boolean,default:!1},orientationMargin:[String,Number]}},P=_({compatConfig:{MODE:3},name:"ADivider",props:b(),setup:function(n,g){var l=g.slots,u=S("divider",n),o=u.prefixCls,h=u.direction,c=r(function(){return n.orientation==="left"&&n.orientationMargin!=null}),f=r(function(){return n.orientation==="right"&&n.orientationMargin!=null}),m=r(function(){var t,i=n.type,x=n.dashed,M=n.plain,e=o.value;return t={},a(t,e,!0),a(t,"".concat(e,"-").concat(i),!0),a(t,"".concat(e,"-dashed"),!!x),a(t,"".concat(e,"-plain"),!!M),a(t,"".concat(e,"-rtl"),h.value==="rtl"),a(t,"".concat(e,"-no-default-orientation-margin-left"),c.value),a(t,"".concat(e,"-no-default-orientation-margin-right"),f.value),t}),y=r(function(){var t=typeof n.orientationMargin=="number"?"".concat(n.orientationMargin,"px"):n.orientationMargin;return d(d({},c.value&&{marginLeft:t}),f.value&&{marginRight:t})}),C=r(function(){return n.orientation.length>0?"-"+n.orientation:n.orientation});return function(){var t,i=D((t=l.default)===null||t===void 0?void 0:t.call(l));return v("div",{class:[m.value,i.length?"".concat(o.value,"-with-text ").concat(o.value,"-with-text").concat(C.value):""],role:"separator"},[i.length?v("span",{class:"".concat(o.value,"-inner-text"),style:y.value},[i]):null])}}});const I=p(P);export{I as D};