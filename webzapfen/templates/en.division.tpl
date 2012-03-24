{{define "DivisionTemplate"}}<!DOCTYPE HTML>
<html>
  <head>
    <title>Division</title>
    <style type="text/css">
      .divisionColumn {
	border: 1px solid black;
	display: inline-block;
	float: left;
    }
      .resultColumn {
	display: inline-block;
	float: left;
    }
    </style>
  </head>
  <body>
  <header>
    <h1>Division</h1>
  </header>{{if .Error}}
  <p>{{range .Error}}
    {{.}}<br/>{{end}}
  </p>{{end}}
  <form>
    Dividend: <input type="text" name="dividend" value="{{.Dividend}}" autofocus="true"/>
    Divisor:<input type="text" name="divisor" value="{{.Divisor}}" />
    Precision:<input type="text" name="prec" value="{{.Precision}}" />
    Stop when remainder is zero:<input type="checkbox" name="stopremz" value="true" {{if .StopRemz}}checked="checked"{{end}}/>
    <input type="submit" value="Submit" />
  </form>
  <pre>
{{.IntermediateStr}}
  </pre>
  {{tplfuncdivdisplay .Intermediate}}
  </body>
</html>{{end}}