package main

import (
	"io"
	"text/template"
)

func html_render(w io.Writer, filepath string, html []byte) error {

	tpl := `<!DOCTYPE html>
  <html>
  <head>
  <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
  <title></title>
  <script src="/!script" ></script>
  <script>
  function setSize(){
    var left_w = 300;

    var workground_h = $(window).height();
    var workground_w = $(window).width();

    $('#sidemenu').width(left_w).height(workground_h).offset({top:0, left:0});
}

function rebuild_index(){
    $('#reindex_btn').attr('disabled','disabled');
    $.ajax({
      url: "/!rebuild"
    }).complete(function(){
      $('#reindex_btn').removeAttr('disabled');
    });
}

function set_path_ipt(fpath){
  if($('#current_path_ipt').val()!=fpath){
    $('#current_path_ipt').val(fpath);
  }
}

function update(f_url, target){
    $.ajax({
      url: f_url
    }).success(function(data) {
      $(target).html(data);
      if(target=="#sidemenu-body"){
        $('#sidemenu-body a').each(function(i,a){
            var addr = $(a).attr('href');
            var first_char = String(addr).substr(0,1);
            if(first_char!='#'){
              if(first_char=='/'){
                addr = String(addr).substr(1);
              }
                var nurl = "/" + addr;
                $(a).attr('href', nurl);
                if(nurl==current_path){
                    $(a).addClass("active");
                }
            }
        });
        setTimeout(function(){update(f_url, target)}, 1000)
      }else{
        window.document.title = $('#main h1').html();
        setTimeout(function(){update(current_path+"?one=true", target)}, 1000)
      }
    });
}

var current_path = window.location.pathname;

$(function(){
    window.document.title = $('#main h1').html();
    setSize();
    $(window).bind("resize", setSize);
    update("/SUMMARY.md?one=summary", "#sidemenu-body");
    update(current_path+"?one=true", "#main-body");
    $('#sidemenu-body').bind('click',function(e){
        e.stopImmediatePropagation();
        current_path = $(e.target).attr('href');
        return false;
      })
})
  </script>
<style>
#sidemenu{font-size:0.9em;width:300px;border-right:2px solid #f0f0f0;position:fixed;overflow:scroll;}
#main-body{max-width:700px;}
#sidemenu .active{font-weight:bold;color:#000;border-right:8px solid #900;padding-right:5px}
#main{margin-left:350px;padding:0}
#main .filepath{border:none;width:100%%;color:#009;font-size:1.2em}
#reindex{position:absolute;right:5px;top:5px;padding:0;margin:0}

.doc-style{padding:20px}
@font-face {
  font-family: octicons-anchor;
  src: url(data:font/woff;charset=utf-8;base64,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) format('woff');
}

.doc-style {
  -webkit-text-size-adjust: 100%;
  -ms-text-size-adjust: 100%;
  text-size-adjust: 100%;
  color: #333;
  overflow: hidden;
  font-family: "Helvetica Neue", Helvetica, "Segoe UI", Arial, freesans, sans-serif;
  font-size: 16px;
  line-height: 1.6;
  word-wrap: break-word;
}

.doc-style a {
  background-color: transparent;
}

.doc-style a:active,
.doc-style a:hover {
  outline: 0;
}

.doc-style strong {
  font-weight: bold;
}

.doc-style h1 {
  font-size: 2em;
  margin: 0.67em 0;
}

.doc-style img {
  border: 0;
}

.doc-style hr {
  box-sizing: content-box;
  height: 0;
}

.doc-style pre {
  overflow: auto;
}

.doc-style code,
.doc-style kbd,
.doc-style pre {
  font-family: monospace, monospace;
  font-size: 1em;
}

.doc-style input {
  color: inherit;
  font: inherit;
  margin: 0;
}

.doc-style html input[disabled] {
  cursor: default;
}

.doc-style input {
  line-height: normal;
}

.doc-style input[type="checkbox"] {
  box-sizing: border-box;
  padding: 0;
}

.doc-style table {
  border-collapse: collapse;
  border-spacing: 0;
}

.doc-style td,
.doc-style th {
  padding: 0;
}

.doc-style * {
  box-sizing: border-box;
}

.doc-style input {
  font: 13px/1.4 Helvetica, arial, nimbussansl, liberationsans, freesans, clean, sans-serif, "Segoe UI Emoji", "Segoe UI Symbol";
}

.doc-style a {
  color: #4078c0;
  text-decoration: none;
}

.doc-style a:hover,
.doc-style a:active {
  text-decoration: underline;
}

.doc-style hr {
  height: 0;
  margin: 15px 0;
  overflow: hidden;
  background: transparent;
  border: 0;
  border-bottom: 1px solid #ddd;
}

.doc-style hr:before {
  display: table;
  content: "";
}

.doc-style hr:after {
  display: table;
  clear: both;
  content: "";
}

.doc-style h1,
.doc-style h2,
.doc-style h3,
.doc-style h4,
.doc-style h5,
.doc-style h6 {
  margin-top: 15px;
  margin-bottom: 15px;
  line-height: 1.1;
}

.doc-style h1 {
  font-size: 22px;
}

.doc-style h2 {
  font-size: 18px;
}

.doc-style h3 {
  font-size: 16px;
}

.doc-style h4 {
  font-size: 14px;
}

.doc-style h5 {
  font-size: 12px;
}

.doc-style h6 {
  font-size: 11px;
}

.doc-style blockquote {
  margin: 0;
}

.doc-style ul,
.doc-style ol {
  padding: 0;
  margin-top: 0;
  margin-bottom: 0;
}

.doc-style dd {
  margin-left: 0;
}

.doc-style code {
  font-family: Consolas, "Liberation Mono", Menlo, Courier, monospace;
  font-size: 12px;
}

.doc-style pre {
  margin-top: 0;
  margin-bottom: 0;
  font: 12px Consolas, "Liberation Mono", Menlo, Courier, monospace;
}

.doc-style .octicon {
  font: normal normal normal 16px/1 octicons-anchor;
  display: inline-block;
  text-decoration: none;
  text-rendering: auto;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

.doc-style .octicon-link:before {
  content: '\f05c';
}

.doc-style>*:first-child {
  margin-top: 0 !important;
}

.doc-style>*:last-child {
  margin-bottom: 0 !important;
}

.doc-style a:not([href]) {
  color: inherit;
  text-decoration: none;
}

.doc-style .anchor {
  position: absolute;
  top: 0;
  left: 0;
  display: block;
  padding-right: 6px;
  padding-left: 30px;
  margin-left: -30px;
}

.doc-style .anchor:focus {
  outline: none;
}

.doc-style h1,
.doc-style h2,
.doc-style h3,
.doc-style h4,
.doc-style h5,
.doc-style h6 {
  position: relative;
  margin-top: 1em;
  margin-bottom: 16px;
  font-weight: bold;
  line-height: 1.4;
}

.doc-style h1 .octicon-link,
.doc-style h2 .octicon-link,
.doc-style h3 .octicon-link,
.doc-style h4 .octicon-link,
.doc-style h5 .octicon-link,
.doc-style h6 .octicon-link {
  display: none;
  color: #000;
  vertical-align: middle;
}

.doc-style h1:hover .anchor,
.doc-style h2:hover .anchor,
.doc-style h3:hover .anchor,
.doc-style h4:hover .anchor,
.doc-style h5:hover .anchor,
.doc-style h6:hover .anchor {
  padding-left: 8px;
  margin-left: -30px;
  text-decoration: none;
}

.doc-style h1:hover .anchor .octicon-link,
.doc-style h2:hover .anchor .octicon-link,
.doc-style h3:hover .anchor .octicon-link,
.doc-style h4:hover .anchor .octicon-link,
.doc-style h5:hover .anchor .octicon-link,
.doc-style h6:hover .anchor .octicon-link {
  display: inline-block;
}

.doc-style h1 {
  padding-bottom: 0.3em;
  line-height: 1.2;
  border-bottom: 1px solid #eee;
}

.doc-style h1 .anchor {
  line-height: 1;
}

.doc-style h2 {
  padding-bottom: 0.3em;
  line-height: 1.225;
  border-bottom: 1px solid #eee;
}

.doc-style h2 .anchor {
  line-height: 1;
}

.doc-style h3 {
  line-height: 1.43;
}

.doc-style h3 .anchor {
  line-height: 1.2;
}

.doc-style h4 .anchor {
  line-height: 1.2;
}

.doc-style h5 .anchor {
  line-height: 1.1;
}

.doc-style h6 {
  color: #777;
}

.doc-style h6 .anchor {
  line-height: 1.1;
}

.doc-style p,
.doc-style blockquote,
.doc-style ul,
.doc-style ol,
.doc-style dl,
.doc-style table,
.doc-style pre {
  margin-top: 0;
  margin-bottom: 16px;
}

.doc-style hr {
  height: 4px;
  padding: 0;
  margin: 16px 0;
  background-color: #e7e7e7;
  border: 0 none;
}

.doc-style ul,
.doc-style ol {
  padding-left: 2em;
}

.doc-style ul ul,
.doc-style ul ol,
.doc-style ol ol,
.doc-style ol ul {
  margin-top: 0;
  margin-bottom: 0;
}

.doc-style li>p {
  margin-top: 16px;
}

.doc-style dl {
  padding: 0;
}

.doc-style dl dt {
  padding: 0;
  margin-top: 16px;
  font-size: 1em;
  font-style: italic;
  font-weight: bold;
}

.doc-style dl dd {
  padding: 0 16px;
  margin-bottom: 16px;
}

.doc-style blockquote {
  padding: 0 15px;
  color: #777;
  border-left: 4px solid #ddd;
}

.doc-style blockquote>:first-child {
  margin-top: 0;
}

.doc-style blockquote>:last-child {
  margin-bottom: 0;
}

.doc-style table {
  display: block;
  width: 100%;
  overflow: auto;
  word-break: normal;
  word-break: keep-all;
}

.doc-style table th {
  font-weight: bold;
}

.doc-style table th,
.doc-style table td {
  padding: 6px 13px;
  border: 1px solid #ddd;
}

.doc-style table tr {
  background-color: #fff;
  border-top: 1px solid #ccc;
}

.doc-style table tr:nth-child(2n) {
  background-color: #f8f8f8;
}

.doc-style img {
  max-width: 100%;
  box-sizing: border-box;
}

.doc-style code {
  padding: 0;
  padding-top: 0.2em;
  padding-bottom: 0.2em;
  margin: 0;
  font-size: 85%;
  background-color: rgba(0,0,0,0.04);
  border-radius: 3px;
}

.doc-style code:before,
.doc-style code:after {
  letter-spacing: -0.2em;
  content: "\00a0";
}

.doc-style pre>code {
  padding: 0;
  margin: 0;
  font-size: 100%;
  word-break: normal;
  white-space: pre;
  background: transparent;
  border: 0;
}

.doc-style .highlight {
  margin-bottom: 16px;
}

.doc-style .highlight pre,
.doc-style pre {
  padding: 16px;
  overflow: auto;
  font-size: 85%;
  line-height: 1.45;
  background-color: #f7f7f7;
  border-radius: 3px;
}

.doc-style .highlight pre {
  margin-bottom: 0;
  word-break: normal;
}

.doc-style pre {
  word-wrap: normal;
}

.doc-style pre code {
  display: inline;
  max-width: initial;
  padding: 0;
  margin: 0;
  overflow: initial;
  line-height: inherit;
  word-wrap: normal;
  background-color: transparent;
  border: 0;
}

.doc-style pre code:before,
.doc-style pre code:after {
  content: normal;
}

.doc-style kbd {
  display: inline-block;
  padding: 3px 5px;
  font-size: 11px;
  line-height: 10px;
  color: #555;
  vertical-align: middle;
  background-color: #fcfcfc;
  border: solid 1px #ccc;
  border-bottom-color: #bbb;
  border-radius: 3px;
  box-shadow: inset 0 -1px 0 #bbb;
}

.doc-style .pl-c {
  color: #969896;
}

.doc-style .pl-c1,
.doc-style .pl-s .pl-v {
  color: #0086b3;
}

.doc-style .pl-e,
.doc-style .pl-en {
  color: #795da3;
}

.doc-style .pl-s .pl-s1,
.doc-style .pl-smi {
  color: #333;
}

.doc-style .pl-ent {
  color: #63a35c;
}

.doc-style .pl-k {
  color: #a71d5d;
}

.doc-style .pl-pds,
.doc-style .pl-s,
.doc-style .pl-s .pl-pse .pl-s1,
.doc-style .pl-sr,
.doc-style .pl-sr .pl-cce,
.doc-style .pl-sr .pl-sra,
.doc-style .pl-sr .pl-sre {
  color: #183691;
}

.doc-style .pl-v {
  color: #ed6a43;
}

.doc-style .pl-id {
  color: #b52a1d;
}

.doc-style .pl-ii {
  background-color: #b52a1d;
  color: #f8f8f8;
}

.doc-style .pl-sr .pl-cce {
  color: #63a35c;
  font-weight: bold;
}

.doc-style .pl-ml {
  color: #693a17;
}

.doc-style .pl-mh,
.doc-style .pl-mh .pl-en,
.doc-style .pl-ms {
  color: #1d3e81;
  font-weight: bold;
}

.doc-style .pl-mq {
  color: #008080;
}

.doc-style .pl-mi {
  color: #333;
  font-style: italic;
}

.doc-style .pl-mb {
  color: #333;
  font-weight: bold;
}

.doc-style .pl-md {
  background-color: #ffecec;
  color: #bd2c00;
}

.doc-style .pl-mi1 {
  background-color: #eaffea;
  color: #55a532;
}

.doc-style .pl-mdr {
  color: #795da3;
  font-weight: bold;
}

.doc-style .pl-mo {
  color: #1d3e81;
}

.doc-style kbd {
  display: inline-block;
  padding: 3px 5px;
  font: 11px Consolas, "Liberation Mono", Menlo, Courier, monospace;
  line-height: 10px;
  color: #555;
  vertical-align: middle;
  background-color: #fcfcfc;
  border: solid 1px #ccc;
  border-bottom-color: #bbb;
  border-radius: 3px;
  box-shadow: inset 0 -1px 0 #bbb;
}

.doc-style .task-list-item {
  list-style-type: none;
}

.doc-style .task-list-item+.task-list-item {
  margin-top: 3px;
}

.doc-style .task-list-item input {
  margin: 0 0.35em 0.25em -1.6em;
  vertical-align: middle;
}

.doc-style :checked+.radio-label {
  z-index: 1;
  position: relative;
  border-color: #4078c0;
}
    </style>
    </head>
    <body>
    <div id="sidemenu">
      <div id="reindex"><button id="reindex_btn" onclick="rebuild_index()">重建索引</button></div>
      <div class="doc-style" id="sidemenu-body"></div>
    </div>
    <div id="main">
      <input id="current_path_ipt" class="filepath" value="{{.filepath}}" onclick="this.select();" />
      <div class="doc-style" id="main-body">
      {{.body}}
      </div>
    </div>
    </body>
    </html>`

	v := map[string]interface{}{
		"body":     string(html),
		"filepath": filepath,
	}

	t, err := template.New("foo").Parse(tpl)
	if err != nil {
		panic(err)
	}
	return t.Execute(w, v)
}
