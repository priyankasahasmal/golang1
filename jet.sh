!/bin/bash

go run github.com/go-jet/jet/cmd/jet \
  -source=postgresql \
  -host=dpg-d1l9vbbe5dus73feni1g-a.oregon-postgres.render.com \
  -port=5432 \
  -user=database1_mi2h_user \
  -password=jqy7Bg3KKql230L8LMLJfArda5ziDWPS \
  -dbname=database1_mi2h \
  -schema=public \
  -path=./generated \
  -sslmode=require 
