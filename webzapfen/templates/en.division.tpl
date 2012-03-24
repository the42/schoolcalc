{{define "DivisionTemplate"}}<!DOCTYPE HTML>
<html>
  <head>
    <title>Division</title>
    <style type="text/css">
      .divisionColumn[division=true][boxed=true] {
	border: 1px solid black;
	display: inline-block;
	float: left;
    }
      .divisionColumn[division=true] {
	display: inline-block;
	float: left;
    }
      .divisionColumn[result=true] {
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
    <p>
      Dividend: <input type="text" name="dividend" value="{{.Dividend}}" autofocus="true"/>
      Divisor: <input type="text" name="divisor" value="{{.Divisor}}" />
      <input type="submit" value="Submit" />
    </p>
    <p>
      Precision: <input type="text" name="prec" value="{{.Precision}}" />
      Continue calculating until precision, even if remainder is already zero <input type="checkbox" name="stopremz" value="false"{{if not .StopRemz}} checked="checked"{{end}}/>
      Display boxes are off <input type="checkbox" name="boxed" value="false"{{if not .Boxed}} checked="checked"{{end}}/>
    </p>
  </form>
  <pre>
{{.IntermediateStr}}
  </pre>
  {{tplfuncdivdisplay .Intermediate .Boxed}}
  </body>
</html>{{end}}