{{define "Title"}}Excersises{{end}}{{define "PathToStatic"}}../static/{{end}}{{define "Payload"}}
  <header>
    <h1>Excersises</h1>
  </header>
  <form>
    <p>
      <label for="levelsetter">Difficulty level</label>
      <select name="level" id="levelsetter">
	{{setIntOptionSelected 0 .Level "Select excercise level..."}}
	{{setIntOptionSelected 6 .Level "user defined"}}
	{{setIntOptionSelected 1 .Level "Beginner"}}
	{{setIntOptionSelected 2 .Level "Apprentice"}}
	{{setIntOptionSelected 3 .Level "Sophomore"}}
	{{setIntOptionSelected 4 .Level "Advanced"}}
	{{setIntOptionSelected 5 .Level "Master"}}
      </select>
    </p>
    <p>
      <label for="n">Number of Excersises:</label>
      <input type="text" name="n" value="{{.NumberofExcersises}}" id="n" pattern="\d+" title="a number" size="2" autofocus="true"/> 
    </p>
    <p>
      <input type="submit" value="Submit"/>
    </p>
    <p>
      <button id="visibleb" type="button"></button>
      <input type="hidden" id="visiblef" name="visiblef" value="{{.Visibility}}"/>
    </p>
    <div id="excercisedetails">
      <section class="group1">
        <div class="inputArea">
	  <h3>Dividend/Divisor Size</h3>
	  <p>
	    <label for="dividendrange">Range of digits in dividend</label>
	    <input type="text" class="excercisedetail" id="dividendrange" name="dividendrange" pattern="\d+(-\d+)?" title="a number eg. 2 which means from zero to 2 or a range, eg. 2-3, which means 2 to 3 digits" value="{{.DividendRange}}" size="6"/>
	    <label for="signdividend">Sign of dividend</label>
	    <select name="signdividend" id="signdividend" class="changeexcercisedetail">
	      {{setIntOptionSelected -1 .SignDividend "positive"}}
	      {{setIntOptionSelected 1 .SignDividend "negative"}}
	      {{setIntOptionSelected 0 .SignDividend "automatic"}}
	    </select>
	  </p>
	  <p>
	    <label for="divisorrange">Range of digits in divisor</label>
	    <input type="text" class="excercisedetail" id="divisorrange" name="divisorrange" pattern="\d+(-\d+)?" title="a number eg. 2 which means from zero to 2 or a range, eg. 2-3, which means 2 to 3 digits" value="{{.DivisorRange}}" size="6"/>
	    <label for="signdivisor">Sign of divisor</label>
	    <select name="signdivisor" id="signdivisor" class="changeexcercisedetail">
	      {{setIntOptionSelected -1 .SignDivisor "positive"}}
	      {{setIntOptionSelected 1 .SignDivisor "negative"}}
	      {{setIntOptionSelected 0 .SignDivisor "automatic"}}
	    </select>
	  </p>
        </div>
      </section>
      <section class="group2">
        <div class="inputArea">
	  <h3>Digit Range</h3>
	  <p>
	    <label for="divisornumrange">Rounded divisor consists of</label>
	    <input type="text" name="divisornumrange" id="divisornumrange" pattern="\d(,\d)*" title="digits, separated by comma, eg. 1,3" list="predivisornumrange" value="{{.DivisorNumRange}}"/>
	    <datalist id="predivisornumrange">
	      <option value="0,2,5">2, 5 and zero</option>
	      <option value="0,2,5,4,3,6,8">2, 5, 4, 3, 6, 8 and zero</option>
	      <option value="0,1,2,3,4,5,6,7,8,9">all digits</option>
	      <option>enter digits ...</option>
	    </datalist>
	  </p>
	  <p>
	    <label for="dividendnumrange">Dividend contains</label>
	    <input type="text" name="dividendnumrange" id="dividendnumrange" pattern="\d(,\d)*" title="digits, separated by comma, eg. 1,3" list="predividendnumrange" value="{{.DividendNumRange}}"/>
	    <datalist id="predividendnumrange">
	      <option value="0,1,2,3,4,5">the digits upto 5 and zero</option>
	      <option value="6,7,8,9,0">the digits greater than 5 and zero</option>
	      <option value="0,1,2,3,4,5,6,7,8,9">all digits</option>
	      <option>enter digits ...</option>
	    </datalist>
	  </p>
        </div>
      </section>
      <section class="group3">
        <div class="inputArea">
	  <h3>Decimal Places</h3>
	  <p>
	    <label for="numremz">Number of decimal places</label>
	    <input type="text" class="excercisedetail" name="numremz" id="numremz" pattern="\d*" title="only numbers allowed" value="{{.MaxDigitisPastPointUntilZero}}" size="10"/> 
	  </p>
        </div>
      </section>
    </div>
  </form>
  <script type="text/javascript">
    function setInputFieldstoLevel() {
      switch(document.forms[0].level.value) {
      case "1":
        document.forms[0].dividendrange.value = "1-2";
        break;
      case "5":
        alert("You rock!");
        break;
      }
    }

    function setSelectiontoDefault() {
      document.forms[0].level[0].selected = "6";
    }

    function alterVisibility() {
      var area = document.getElementById('excercisedetails');
      var buttonb = document.getElementById('visibleb');
      var visiblef = document.getElementById('visiblef');
      if ( getComputedStyle(area).getPropertyValue('display') == 'none') {
        area.style.display = 'block';
	buttonb.innerHTML="Hide details"
	visiblef.value="on"
      } else {
        area.style.display = 'none';
	buttonb.innerHTML="Show details"
        visiblef.value=""
      }
    }

    var area = document.getElementById('excercisedetails');
    var buttonb = document.getElementById('visibleb');
    if ( document.forms[0].visiblef.value=="on") {
      buttonb.innerHTML="Hide details"
      area.style.display = 'block';
    } else {
      area.style.display = 'none';
      buttonb.innerHTML="Show details"
    }  

    var buttonvisiblearea = document.getElementById('visibleb');
    buttonvisiblearea.addEventListener('click', alterVisibility);

    var details = document.getElementsByClassName('excercisedetail');
    for(i=0; i < details.length; i++) {
	details[i].addEventListener('input',  setSelectiontoDefault);
    }
    var details = document.getElementsByClassName('changeexcercisedetail');
    for(i=0; i < details.length; i++) {
	details[i].addEventListener('change',  setSelectiontoDefault);
    }

    document.getElementById('levelsetter').addEventListener('change', setInputFieldstoLevel);
  </script>{{end}}