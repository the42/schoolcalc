{{define "Title"}}Excersises{{end}}{{define "PathToStatic"}}../static/{{end}}{{define "Payload"}}
  <header>
    <h1>Excersises</h1>
  </header>
  <form>
    <p>
      <select name="level">
        {{setOptionSelected "0" .Level "Select excercise level..."}}
        {{setOptionSelected "1" .Level "Beginner"}}
        {{setOptionSelected "2" .Level "Apprentice"}}
        {{setOptionSelected "3" .Level "Sophomore"}}
        {{setOptionSelected "4" .Level "Advanced"}}
        {{setOptionSelected "5" .Level "Master"}}
      </select>
    </p>
    <p>
      Number of Excersises:  <input type="text" name="n" value="{{.NumberofExcersises}}" autofocus="true"/> 
    </p>
    <p>
      Minimum Dividend: <input type="text" name="mindividend" value="{{.MinDividend}}" />
      Maximum Dividend: <input type="text" name="maxdividend" value="{{.MaxDividend}}"/>
    </p>
    <p>
      Minimum Divisor: <input type="text" name="mindivisor" value="{{.MinDivisor}}" />
      Maximum Dividend: <input type="text" name="maxdivisor" value="{{.MaxDivisor}}" />
    </p>
    <p>
      Number of Remainders past point:  <input type="text" name="numremz" value="{{.MaxDigitisPastPointUntilZero}}" /> 
    </p>
    <input type="submit" value="Submit" />
  </form>{{end}}