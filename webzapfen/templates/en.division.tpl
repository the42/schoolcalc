<!DOCTYPE HTML>
<html>
  <head>
    <title>Division</title>
  </head>
  <body>
  <header>
    <h1>Division</h1>
  </header>{{if .Error}}
  <p>{{range .Error}}
    {{.}}<br/>{{end}}
  </p>{{end}}
  <form>
    Dividend: <input type="text" name="dividend" value="{{.Dividend}}" />
    Divisor:<input type="text" name="divisor" value="{{.Divisor}}" />
    <input type="submit" value="Submit" />
  </form> 
  <pre>{{.Intermediate}}</pre>
  </body>
</html>