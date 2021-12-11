vegeta attack -duration=1s -rate=1000000 -workers=8 -targets="test.http" -output=results.bin
vegeta plot -title=Results .\results > results-plot.html

