{{define "SchoolCalcRoot"}}<!DOCTYPE HTML>
<html>
  <head>
    <!-- style type="text/css">
    .divisionOutputArea {
      background-image:url(data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAACYAAAAmAQMAAACS83vtAAAAAXNSR0IArs4c6QAAAAZQTFRF////mcz/7U0nzgAAAAlwSFlzAAALEwAACxMBAJqcGAAAAAd0SU1FB9wECAceA9EbgDsAAAAYSURBVAjXY2BgUGBgYaAP+R8I/jDQ3UYApJwPAeJX0y0AAAAASUVORK5CYII=);
    }
    .divisionColumn[data-division=true],.divisionColumn[data-result=true] {
      display: inline-block;
    }
    .divisionColumn[data-division=true][data-boxed=true] {
      border: 1px solid black;
    }
    </style -->
    <title>School Division: Divsion and Excercises for Dividing the Pen and Paper Method</title>
  </head>
  <body>
  <header>
    <h1>School Division: Divsion and Excercises for Dividing the Pen and Paper Method</h1>
  </header>
  <p>
    This little tool will help you to learn divide the pen and paper method. Dividing the pen and paper method
    is the hardest beside the three other primary calculating methods of Addition, Subtraction and Multiplication.
  </p>
  <div class="langsel">
  </div>
  <menu>
    What do you want to do?

    Learn
      Learn more about the history of Division
      Learn the Algorithm behind Dividing the pen and paper method

    Test
    Check    
  </menu>
    <!-- Dividend: <input type="text" name="dividend" value="{{.Dividend}}" autofocus="true"/>
    Divisor: <input type="text" name="divisor" value="{{.Divisor}}" />
    <input type="submit" value="Submit" />
  </p>
  <p>
    Precision: <input type="text" name="prec" value="{{.Precision}}" />
    Continue calculating until precision, even if remainder is already zero <input type="checkbox" id="submitprec" name="stopremz" value="false"{{if not .StopRemz}} checked="checked"{{end}}/>
    Display boxes are off <input type="checkbox" id="boxedresult" name="boxed" value="false"{{if not .Boxed}} checked="checked"{{end}}/ -->

  </body>
</html>{{end}}