{{define "DivisionTemplate"}}<!DOCTYPE HTML>
<html>
  <head>
    <style type="text/css">
      .divisionColumn[data-division=true][data-boxed=true] {
	border: 1px solid black;
	display: inline-block;
	float: left;
    }
      .divisionColumn[data-division=true] {
	display: inline-block;
	float: left;
    }
      .divisionColumn[data-result=true] {
	display: inline-block;
	float: left;
    }
    </style>
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
    <p>
      Dividend: <input type="text" name="dividend" value="{{.Dividend}}" autofocus="true"/>
      Divisor: <input type="text" name="divisor" value="{{.Divisor}}" />
      <input type="submit" value="Submit" />
    </p>
    <p>
      Precision: <input type="text" name="prec" value="{{.Precision}}" />
      Continue calculating until precision, even if remainder is already zero <input type="checkbox" id="submitprec" name="stopremz" value="false"{{if not .StopRemz}} checked="checked"{{end}}/>
      Display boxes are off <input type="checkbox" id="boxedresult" name="boxed" value="false"{{if not .Boxed}} checked="checked"{{end}}/>
    </p>
  </form>
  <div class="divisionOutputArea">
  {{tplfuncdivdisplay .SDivide .Boxed}}
  </div>
  <script type="text/javascript">
    function changeBox() {
      var items = document.getElementsByClassName('divisionColumn');
      for(i=0; i < items.length; i++) {
	if(items[i].hasAttribute('data-division')) {
	  items[i].setAttribute('data-boxed', items[i].getAttribute('data-boxed') == 'true' ? 'false' : 'true');
	}
      }
    }
    document.getElementById('boxedresult').addEventListener('click', changeBox);
    document.getElementById('submitprec').addEventListener('click', function() { document.forms[0].submit();} );
  </script>
  </body>
</html>{{end}}