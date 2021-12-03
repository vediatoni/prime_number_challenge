vegeta attack -duration=10s -rate=100 -workers=10 -targets="test.http" -output=results.bin
vegeta plot -title=Results .\results > results-plot.html