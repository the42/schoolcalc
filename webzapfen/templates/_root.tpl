{{define "DivisionSetup"}}<!DOCTYPE HTML>
<html lang="{{.CurrLang}}">
  <head>
  <meta charset="UTF-8">
    <title>{{template "Title"}}</title>
    <link rel="shortcut icon" href="{{template "PathToStatic"}}images/favicon.png" type="image/png"/> 
    <link rel="icon" href="{{template "PathToStatic"}}images/favicon.png" type="image/png"/>
    <link rel="stylesheet" type="text/css" href="{{template "PathToStatic"}}css/styles.css"/>
  </head>
  <body>
  <div class="bg">
    <!--start container-->
    <div id="container">
    <!--start header-->
    <header>
      <!--start logo-->
      <a href="#" id="logo"><img src="static/images/logo.png" width="180" height="43" alt="logo"/></a>    
      <!--end logo-->
      <!--start menu-->
      <nav>
        <ul>
          <li><a href="#" class="current">Home</a></li>
          <li><a href="#">News</a></li>
       </ul>{{langselector .CurrLang}}
      </nav>
     <!--end menu-->
    </header>{{if .Error}}
    <p>{{range .Error}}
    {{.}}<br/>{{end}}
    </p>{{end}}{{template "Payload" .}}
    </div>
    <footer>
      <div class="container">  
        <div id="FooterTwo"> © 2011 Minimalism </div>
        <div id="FooterTree"> Valid html5, design and code by <a href="http://www.marijazaric.com">marija zaric - creative simplicity</a> </div> 
      </div>
   </footer>
   </div>
   <!--end bg-->
  </body>
</html>{{end}}