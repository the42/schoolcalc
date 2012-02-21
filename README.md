schoolcalc - Support pupils when learning to divide
===================================================

Short description
-----------------
A library for

* dividing integers or rationals. It calculates intermidate steps when performing divisions.
Supports negative numerals.
* Zapfenrechnung. A Zapfenrechnung (Zapfen rechnen) looks like this:

<pre><code>
         27 * 2 = 54
         54 * 3 = 162
        162 * 4 = 648
        648 * 5 = 3240
       3240 * 6 = 19440
      19440 * 7 = 136080
     136080 * 8 = 1088640
    1088640 * 9 = 9797760
    9797760 / 2 = 4898880
    4898880 / 3 = 1632960
    1632960 / 4 = 408240
     408240 / 5 = 81648
      81648 / 6 = 13608
      13608 / 7 = 1944
       1944 / 8 = 243
        243 / 9 = 27
</code></pre>

Learning to divide is particulary hard for pupils and this library
provides a method which returns intermediate steps and indention hints
to display those intermediate remainders on an output device.


Features
--------

* Arbitrary precision mathematics
* Division: Negative numbers & fractional numbers supported when deviding the pen & paper method
* Division: Intermediate steps are returned for display / control

Installation
------------

Library:
  go install github.com/the42/schoolcalc

Programs:
  Command line program sdivcon

Usage
-----

All features provided by the package are covered by test cases.

License
-------

The package is released under the [Simplified BSD
License](http://www.freebsd.org/copyright/freebsd-license.html) See file
"LICENSE"

Testing
-------

To run the tests:

  go test github.com/the42/schoolcalc
