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
    a.emptylink {
      color:blue;
      text-decoration: none;
    }
    a.emptylink:hover {
      color:blue;
      font-weight:bold;
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
    Number: <input type="text" name="number" value="{{.Number}}" autofocus="true"/>
    <input type="submit" value="Submit" />
  </form>{{if .Zapfen}}
  <a href='#' id='toggleallintermedsteps' class='emptylink'>Toggle</a> display of intermediate division steps
  <div class="zapfenOutputArea">
    <table>
      {{tplfunczapfendisplay .Zapfen .IntermedZapfen}}
    </table>
  </div>
  <script type="text/javascript">
    function changeallIntermediate() {
      var items = document.getElementsByClassName('zapfenintermeddivisionrow');
      for(i=0; i < items.length; i++) {
	items[i].style.display = (getComputedStyle(items[i]).getPropertyValue('display') == 'none') ? 'table-row' : 'none';
      }
    }

    function changeIntermediate(e) {
      var num = this.attributes['id'].value.match(/\d+/);
      var item = document.getElementById('zapfenintermeddivisionrow' + num);
      item.style.display = (getComputedStyle(item).getPropertyValue('display') == 'none') ? 'table-row' : 'none';
    }

    document.getElementById('toggleallintermedsteps').addEventListener('click', changeallIntermediate);
    var dividends = document.getElementsByClassName('zapfendividenditem');
    for(i=0; i < dividends.length; i++) {
	dividends[i].addEventListener('click',  changeIntermediate);
    }
  </script>{{end}}
  </body>
</html>{{end}}