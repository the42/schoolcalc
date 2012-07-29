{{define "Title"}}Excersises{{end}}{{define "PathToStatic"}}../static/{{end}}{{define "Payload"}}
  <header>
    <h1>Excersises</h1>
  </header>
  <form>
    <div>
      <p>
	<select name="level" id="levelsetter">
	  {{setLevelOptionSelected 0 .Level "Select excercise level..."}}
	  {{setLevelOptionSelected 1 .Level "Beginner"}}
	  {{setLevelOptionSelected 2 .Level "Apprentice"}}
	  {{setLevelOptionSelected 3 .Level "Sophomore"}}
	  {{setLevelOptionSelected 4 .Level "Advanced"}}
	  {{setLevelOptionSelected 5 .Level "Master"}}
	</select>
      </p>
      <p>
	Number of Excersises:  <input type="text" name="n" value="{{.NumberofExcersises}}" autofocus="true"/> 
      </p>
      <p>
	<a href="#" id='togglevisibility' ></a>
      </p>
    </div>
    <div id="excercisedetails">
      <p>
	Minimum Dividend: <input type="text" class="excercisedetail" name="mindividend" value="{{.MinDividend}}" />
	Maximum Dividend: <input type="text" class="excercisedetail" name="maxdividend" value="{{.MaxDividend}}"/>
      </p>
      <p>
	Minimum Divisor: <input type="text" class="excercisedetail" name="mindivisor" value="{{.MinDivisor}}" />
	Maximum Dividend: <input type="text" class="excercisedetail" name="maxdivisor" value="{{.MaxDivisor}}" />
      </p>
      <p>
	Number of Remainders past point:  <input type="text" class="excercisedetail" name="numremz" value="{{.MaxDigitisPastPointUntilZero}}" /> 
      </p>
    </div>
    <input type="submit" value="Submit" />
  </form>
  <script type="text/javascript">
    function setInputFieldstoLevel() {
      switch(document.forms[0].level.value) {
      case "1":
        document.forms[0].mindividend.value = "45";
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
    document.getElementById('levelsetter').addEventListener('change', setInputFieldstoLevel);
  </script>{{end}}