{{define "Title"}}Excersises{{end}}{{define "PathToStatic"}}../static/{{end}}{{define "Payload"}}
  <header>
    <h1>Excersises</h1>
  </header>
  <form>
    <p>
      <label for="levelsetter">Difficulty level<label>
      <select name="level" id="levelsetter">
	{{setIntOptionSelected 0 .Level "Select excercise level..."}}
	{{setIntOptionSelected 1 .Level "Beginner"}}
	{{setIntOptionSelected 2 .Level "Apprentice"}}
	{{setIntOptionSelected 3 .Level "Sophomore"}}
	{{setIntOptionSelected 4 .Level "Advanced"}}
	{{setIntOptionSelected 5 .Level "Master"}}
      </select>
    </p>
    <p>
      <label for="n">Number of Excersises:</label>
      <input type="text" name="n" value="{{.NumberofExcersises}}" id="n" size="2" autofocus="true"/> 
    </p>
    <p>
      <input type="submit" value="Submit"/>
    </p>
    <p>
      <a href="#" id="togglevisibility" ></a>
    </p>
    <div id="excercisedetails">
      <section class="group1">
        <div class="inputArea">
	  <h3>Dividend/Divisor Size</h3>
	  <p>
	    <label for="dividendrange">Range of digits in dividend</label>
	    <input type="text" class="excercisedetail" id="dividendrange" name="dividendrange" value="{{.DividendRange}}" size="6"/>
	    <label for="signdividend">Sign of dividend</label>
	    <select name="signdividend" id="signdividend" class="changeexcercisedetail">
	      {{setIntOptionSelected -1 .SignDividend "positive"}}
	      {{setIntOptionSelected 1 .SignDividend "negative"}}
	      {{setIntOptionSelected 0 .SignDividend "automatic"}}
	    </select>
	  </p>
	  <p>
	    <label for="divisorrange">Range of digits in divisor</label>
	    <input type="text" class="excercisedetail" id="divisorrange" name="divisorrange" value="{{.DivisorRange}}" size="6"/>
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
	    <label for="divisornumrange">Rounded divisor contains</label>
	    <input type="text" class="excercisedetail" id="divisornumrange" name="divisornumrange" value="{{.DivisorNumRange}}" size="6"/>
	  </p>
	  <p>
	    <label for="dividendnumrange">Dividend contains</label>
	    <input type="text" class="excercisedetail" id="dividendnumrange" name="dividendnumrange" value="{{.DividendNumRange}}" size="6"/>
	  </p>
        </div>
      </section>
      <section class="group3">
        <div class="inputArea">
	  <h3>Decimal Places</h3>
	  <p>
	    <label for="numremz">Number of decimal places</label>
	    <input type="text" class="excercisedetail" name="numremz" id="numremz" value="{{.MaxDigitisPastPointUntilZero}}" size="10"/> 
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
      document.forms[0].level[0].selected = "1";
    }

    function alterVisibility() {
      var area = document.getElementById('excercisedetails');
      var button = document.getElementById('togglevisibility');
      if ( getComputedStyle(area).getPropertyValue('display') == 'none') {
        area.style.display = 'block';
        button.innerHTML = "Hide details";
      } else {
        area.style.display = 'none';
        button.innerHTML = "Show details";
      }
    }

    var buttonvisiblearea = document.getElementById('togglevisibility');
    buttonvisiblearea.innerHTML="Hide details";
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