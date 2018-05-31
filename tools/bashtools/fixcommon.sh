#!/bin/bash

cp g2.csv t.csv

perl -pe 's/gmail"/gmail\.com"/'                   t.csv > t1.csv ; mv t1.csv t.csv
perl -pe 's/([^@])gmail\.com"/$1\@gmail.com"/'     t.csv > t1.csv ; mv t1.csv t.csv
perl -pe 's/yahoo"/yahoo\.com"/'                   t.csv > t1.csv ; mv t1.csv t.csv
perl -pe 's/([^@])yahoo\.com"/$1\@yahoo.com"/'     t.csv > t1.csv ; mv t1.csv t.csv
perl -pe 's/msn"/msn\.com"/'                       t.csv > t1.csv ; mv t1.csv t.csv
perl -pe 's/([^@])msn\.com"/$1\@msn.com"/'         t.csv > t1.csv ; mv t1.csv t.csv
perl -pe 's/verizon"/verizon\.com"/'               t.csv > t1.csv ; mv t1.csv t.csv
perl -pe 's/([^@])verizon\.com"/$1\@verizon.com"/' t.csv > t1.csv ; mv t1.csv t.csv

perl -pe 's/([^@])yahoocom"/$1\@yahoo.com"/'       t.csv > t1.csv ; mv t1.csv t.csv
mv t.csv g3.csv