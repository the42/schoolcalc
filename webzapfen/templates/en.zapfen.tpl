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
        display:none;
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
  Intermediate division steps are hidden <input type="checkbox" id="toggleallintermedsteps" checked="checked"/>
  <div class="zapfenOutputArea">
    <table>
      {{tplfunczapfendisplay .Zapfen .IntermedZapfen}}
    </table>
  </div>{{end}}
  <script type="text/javascript">
    function changeIntermediate() {
      var items = document.getElementsByClassName('zapfenintermeddivisionrow');
      for(i=0; i < items.length; i++) {
	items[i].style.display = (getComputedStyle(items[i]).getPropertyValue('display') == 'none') ? 'table-row' : 'none';
      }
    }
    document.getElementById('toggleallintermedsteps').addEventListener('click', changeIntermediate);
  </script>
  </body>
</html>{{end}}