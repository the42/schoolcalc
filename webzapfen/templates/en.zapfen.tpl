<!DOCTYPE HTML>
<html>
  <head>
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
  </form> 
  <pre>{{.Intermediate}}</pre>
  </body>
</html>