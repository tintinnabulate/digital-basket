digital-basket
==============

[![Build Status](https://travis-ci.org/tintinnabulate/digital-basket.svg?branch=master)](https://travis-ci.org/tintinnabulate/digital-basket)

Website for taking fixed-price donations.

Notes
-----

Stripe takes a per-transaction fee of 1.4% + 20p (at the time of this writing, 2018-12-03).

The payee must enter their email address and card details every time they wish to make a contribution, *unless* they tick the box "Remember me". If the payee ticks the box "Remember me", a cookie is stored on their device by Stripe which is used to recall their card details, so they need only click "Pay" to pay.

PCI compliance
--------------

The site is compliant to standard PCI DSS v3.2.1. Here's how:

* We use Stripe Checkout to host all form inputs containing card data within an IFRAME served from Stripe’s domain — not ours — so the payee's card information never touches our server.

* We are processing less than 6 million transactions per year, and so are eligible to use a SAQ-A (Self-Assessment Questionnaire A) to prove PCI compliance.

* Our Payment page makes use of a modern version of TLS (TLS 1.2) which significantly reduces the risk of us or our payees being exposed to a man-in-the-middle attack. (Achieved grade A using 'Qualys SSL Server Test': https://www.ssllabs.com/ssltest/index.html)

Traditions compliance
---------------------

The site adheres to Tradition 7. Here's how:

* The homepage includes the statement "In keeping with A.A.’s Seventh Tradition of self-support, we accept contributions only from A.A. members."

The site adheres to Tradition 11. Here's how:

* The charge appears on a payee's bank statement with an anonymised statement descriptor.
