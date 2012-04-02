{{define "ZapfenTemplate"}}<!DOCTYPE HTML>
<html>
  <head>
    <style type="text/css">
      TD.zapfenmultiplier, TD.zapfendividend {
        text-align: right
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
      {{tplfunczapfendisplay .Zapfen}}
    </table>
  </div>{{end}}
  </body>
</html>{{end}}