{{define "ZapfenTemplate"}}<!DOCTYPE HTML>
<html>
  <head>
    <style type="text/css">
      .zapfenmultiplier, .zapfendividend {
        text-align: right;
    }
      .zapfendividendintermed {
        float:right;
    }
      .zapfenintermeddivisionrow {
        vertical-align:top;
    }
      .divisionColumn[data-division=true][data-boxed=true] {
	border: 1px solid black;
	display: inline-block;
	float: left;
    }
    </style>
    <title>Webzapfen</title>
  </head>
  <body>
  <header>
    <h1>Webzapfen</h1>
  </header>{{if .Error}}
  <p>{{range .Error}}
    {{.}}<br/>{{end}}
  </p>{{end}}
  <form>
    Number: <input type="text" name="number" value="{{.Number}}" />
    <input type="submit" value="Submit" />
  </form>{{if .Zapfen}}
  <div class="zapfenOutputArea">
    <table>
      {{tplfunczapfendisplay .Zapfen .IntermedZapfen true}}
    </table>
  </div>{{end}}
  </body>
</html>{{end}}