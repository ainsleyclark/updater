(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-6e527c09"],{"08c7":function(t,e,a){"use strict";a.r(e);var r=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("section",[a("div",{staticClass:"auth-container"},[a("div",{staticClass:"row"},[a("div",{staticClass:"col-12"},[a("header",{staticClass:"header"},[a("div",{staticClass:"header-title"},[a("h1",[t._v("Fields")]),a("Breadcrumbs")],1)])])]),a("div",{staticClass:"row"},[a("div",{staticClass:"col-12"},[a("Alert",{attrs:{colour:"orange"}},[t._t("default",[a("h3",[t._v("Coming soon!")]),a("p",[t._v("Hold tight, we're working on this..")])])],2)],1)])])])},s=[],n=a("90f8"),i=a("a71d"),c={name:"Fields",components:{Breadcrumbs:n["a"],Alert:i["a"]},data:function(){return{}},methods:{}},u=c,o=(a("4ca6"),a("2877")),l=Object(o["a"])(u,r,s,!1,null,"37696838",null);e["default"]=l.exports},"304e":function(t,e,a){},"4c2d":function(t,e,a){"use strict";var r=a("cdcc"),s=a.n(r);s.a},"4ca6":function(t,e,a){"use strict";var r=a("304e"),s=a.n(r);s.a},"5ee2":function(t,e,a){"use strict";var r=a("f836"),s=a.n(r);s.a},"7db0":function(t,e,a){"use strict";var r=a("23e7"),s=a("b727").find,n=a("44d2"),i=a("ae40"),c="find",u=!0,o=i(c);c in[]&&Array(1)[c]((function(){u=!1})),r({target:"Array",proto:!0,forced:u||!o},{find:function(t){return s(this,t,arguments.length>1?arguments[1]:void 0)}}),n(c)},"90f8":function(t,e,a){"use strict";var r=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("nav",{staticClass:"breadcrumbs"},[a("ul",{staticClass:"breadcrumbs-list"},t._l(t.breadcrumbs,(function(e){return a("li",{key:e.url,staticClass:"breadcrumbs-item"},[a("router-link",{staticClass:"breadcrumbs-link",class:{"breadcrumbs-link-active":e.active},attrs:{to:e.url}},[t._v(t._s(e.name)+" ")])],1)})),0)])},s=[],n=(a("7db0"),a("4160"),a("b0c0"),a("ac1f"),a("5319"),a("1276"),a("159b"),{name:"Breadcrumbs",data:function(){return{breadcrumbs:[]}},beforeMount:function(){this.updateList()},watch:{$route:function(){this.updateList()}},methods:{updateList:function(){var t=this;this.breadcrumbs=[];var e=this.$route.fullPath,a=e.split("/");if("home"===this.$route.name)this.breadcrumbs.push({name:"Home",url:"/",active:!0});else{var r="",s=!0;a.forEach((function(e){var a;r+=e+"/",s?(a=r,s=!1):a=r.replace(/\/$/,""),e=e.split("?")[0],"resources"===e&&"settings"===e||t.breadcrumbs.push({name:""===e?"Home":t.capitalize(e),url:a,active:t.$route.fullPath===a})}))}},addPage:function(){},getRoute:function(t){return this.$router.options.routes.find((function(e){return e.name===t}))},capitalize:function(t){return t.replace(/(?:^|\s|["'([{])+\S/g,(function(t){return t.toUpperCase()}))}}}),i=n,c=(a("5ee2"),a("2877")),u=Object(c["a"])(i,r,s,!1,null,"7d86a761",null);e["a"]=u.exports},a71d:function(t,e,a){"use strict";var r=function(){var t=this,e=t.$createElement,a=t._self._c||e;return a("transition",{attrs:{name:"trans-fade-quick",mode:"out-in"}},[t.show?a("div",{staticClass:"alert alert-background",class:"alert-"+t.colour},[a("div",{staticClass:"alert-icon"},["error"===t.type?a("i",{staticClass:"feather feather-alert-triangle"}):t._e(),"warning"===t.type?a("i",{staticClass:"feather feather-alert-circle"}):t._e(),"success"===t.type?a("i",{staticClass:"feather feather-check-circle"}):t._e()]),a("div",{staticClass:"alert-text"},[t._t("default")],2),a("button",{staticClass:"alert-close",attrs:{type:"button","aria-label":"Close"},on:{click:function(e){t.show=!1}}},[a("i",{staticClass:"feather feather-x"})])]):t._e()])},s=[],n={name:"Alert",props:{colour:{type:String,default:""},type:{type:String,default:"error"}},data:function(){return{show:!0}}},i=n,c=(a("4c2d"),a("2877")),u=Object(c["a"])(i,r,s,!1,null,null,null);e["a"]=u.exports},cdcc:function(t,e,a){},f836:function(t,e,a){}}]);
//# sourceMappingURL=chunk-6e527c09.321b1a0e.js.map