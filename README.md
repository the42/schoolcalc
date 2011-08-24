sdivision - Support pupils when learning to divide
==================================================

Short description
-----------------
A library for dividing integers or rationals.
It calculates intermidate steps when performing divisions.
Supports negative numerals.

Learning to divide is particulary hard for pupils and this library
provides a method which returns intermediate steps and indention hints
to display those intermediate remainders on an output device.


Features
--------

* Arbitrary precision mathematics
* Negative numbers & fractional numbers supported
* Intermediate steps are returned for display / control

Installation
------------

Library:
  goinstall github.com/the42/sdivision

Programs:
  More to come. I plan a command line program and probalby sthg. to test the "new" templating

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

  make test
