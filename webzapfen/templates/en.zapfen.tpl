{{define "ZapfenTemplate"}}<!DOCTYPE HTML>
<html lang="{{.CurrLang}}">
  <head>
    <link rel="stylesheet" type="text/css" href="../static/css/styles.css"/>
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
  <button id='toggleallintermedsteps' >Toggle</button> display of intermediate division steps
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